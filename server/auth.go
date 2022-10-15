package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var (
	// clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	// clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
	clientID     = "web"
	clientSecret = "secret"
	issuer       = "http://localhost:9998"
)

// ContextKey is used to store the auth context value in the request context
type ContextKey string

var contextKey ContextKey = "AUTH_CONTEXT"

func newAuthProvider(ctx context.Context) (*authProvider, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("creating OIDC provider: %w", err)
	}
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:3000/api/v1/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "openid", "profile", "preferred_username", "email", "groups"},
	}

	sessioner := newSessioner()
	police := NewPolicyEngine()

	return &authProvider{
		ctx:       ctx,
		oidc:      provider,
		verifier:  verifier,
		config:    config,
		sessioner: sessioner,
		police:    police,
	}, nil
}

type authProvider struct {
	ctx       context.Context
	oidc      *oidc.Provider
	verifier  *oidc.IDTokenVerifier
	config    *oauth2.Config
	sessioner *Sessioner
	police    *PolicyEngine
}

func (p authProvider) authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := p.sessioner.AuthorizeSession(r)
		if err != nil {
			http.Error(w, "no session: "+err.Error(), http.StatusUnauthorized)
			return
		}
		// allow, err := p.police.EvaluateLoginRequest(session.UserInfo)
		// if err != nil {
		// 	http.Error(w, "Error enforcing login policies: "+err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// if !allow {
		// 	http.Error(w, "Forbidden access", http.StatusForbidden)
		// }

		ctx := context.WithValue(r.Context(), contextKey, session)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (p authProvider) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/authenticate", p.handleAuthenticate)
	r.Get("/login", p.handleLogin)
	r.Get("/logout", p.handleLogout)

	r.Get("/auth/callback", p.handleAuthCallback)

	return r
}

func (p authProvider) handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	session, err := p.sessioner.AuthorizeSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	data := struct {
		UserInfo UserInfo `json:"user"`
	}{
		UserInfo: session.UserInfo,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Marshalling JSON", http.StatusInternalServerError)
		return
	}
	w.Write(dataJSON)
}

func (p authProvider) handleLogin(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	nonce := uuid.New().String()
	setCallbackCookie(w, r, "nonce", nonce)
	setCallbackCookie(w, r, "state", state)
	http.Redirect(w, r, p.config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
}

func (p authProvider) handleLogout(w http.ResponseWriter, r *http.Request) {
	if err := p.sessioner.EndSession(r); err != nil {
		http.Error(w, "Cannot logout: "+err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "http://localhost:5173/", http.StatusSeeOther)
}

func (p authProvider) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state == "" {
		http.Error(w, "State not found", http.StatusBadRequest)
		return
	}
	var (
		expectedState = r.URL.Query().Get("state")
	)
	if expectedState != state {
		http.Error(w, fmt.Sprintf("Expected state %q got %q", expectedState, state), http.StatusBadRequest)
		return
	}

	oauth2Token, err := p.config.Exchange(p.ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}
	fmt.Println("rawIDToken: ", rawIDToken)

	idToken, err := p.verifier.Verify(p.ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	nonce, err := r.Cookie("nonce")
	if err != nil {
		http.Error(w, "nonce not found", http.StatusBadRequest)
		return
	}
	if idToken.Nonce != nonce.Value {
		http.Error(w, "nonce did not match", http.StatusBadRequest)
		return
	}

	userInfo, err := p.oidc.UserInfo(p.ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	p.sessioner.NewSession(w, idToken, userInfo)
	http.Redirect(w, r, "http://localhost:5173/", http.StatusSeeOther)
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
