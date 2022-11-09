package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/verifa/coastline/server/oapi"
	"github.com/verifa/coastline/server/session"
	"github.com/verifa/coastline/store"
	"golang.org/x/oauth2"
)

type AuthConfig struct {
	ClientID     string   `envconfig:"client_id"`
	ClientSecret string   `envconfig:"client_secret"`
	Issuer       string   `envconfig:"issuer"`
	Scopes       []string `envconfig:"scopes"`
	RedirectURI  string   `envconfig:"redirect_uri"`
}

// ContextKey is used to store the auth context value in the request context
type ContextKey string

var contextKey ContextKey = "AUTH_CONTEXT"

var devUser = claimsToUser(&session.UserClaims{
	Sub:    "dev",
	Iss:    "dev",
	Email:  "dev@localhost",
	Name:   "dev",
	Groups: []string{"admin", "dev"},
})

func newAuthProvider(ctx context.Context, store *store.Store, config AuthConfig, devMode bool) (*authProvider, error) {
	police := NewPolicyEngine()

	if devMode {
		return &authProvider{
			ctx:         ctx,
			store:       store,
			police:      police,
			devMode:     devMode,
			redirectURI: "/ui",
		}, nil
	}
	provider, err := oidc.NewProvider(ctx, config.Issuer)
	if err != nil {
		return nil, fmt.Errorf("creating OIDC provider: %w", err)
	}
	oidcConfig := &oidc.Config{
		ClientID: config.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  config.RedirectURI,
		Scopes:       config.Scopes,
	}

	return &authProvider{
		ctx:         ctx,
		oidc:        provider,
		verifier:    verifier,
		config:      oauth2Config,
		store:       store,
		police:      police,
		redirectURI: "/ui",
	}, nil
}

type authProvider struct {
	ctx      context.Context
	oidc     *oidc.Provider
	verifier *oidc.IDTokenVerifier
	config   *oauth2.Config
	store    *store.Store
	// sessioner   *Sessioner
	police      *PolicyEngine
	devMode     bool
	redirectURI string
}

func (p authProvider) authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p.devMode {
			ctx := context.WithValue(r.Context(), contextKey, devUser)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		sessionID, err := getSessionCookie(r)
		if err != nil {
			http.Error(w, "getting session cookie: "+err.Error(), http.StatusUnauthorized)
		}
		user, err := p.store.GetSession(sessionID)
		if err != nil {
			http.Error(w, "invalid session: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), contextKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (p authProvider) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/login", p.handleLogin)
	r.Get("/logout", p.handleLogout)
	r.Get("/auth/callback", p.handleAuthCallback)

	return r
}

func (p authProvider) UserInfo(r *http.Request) (*oapi.User, error) {
	sessionID, err := getSessionCookie(r)
	if err != nil {
		return nil, fmt.Errorf("getting session cookie: %w", err)
	}
	user, err := p.store.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session: %w", err)
	}
	return user, nil
}

func (p authProvider) handleLogin(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	nonce := uuid.New().String()
	if p.devMode {
		sessionID, err := p.store.NewSession(devUser)
		if err != nil {
			http.Error(w, "Creating new session: "+err.Error(), http.StatusInternalServerError)
		}
		writeSessionCookieHeader(w, sessionID)
		http.Redirect(w, r, p.redirectURI, http.StatusSeeOther)
		return
	}
	setCallbackCookie(w, r, "nonce", nonce)
	setCallbackCookie(w, r, "state", state)
	http.Redirect(w, r, p.config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
}

func (p authProvider) handleLogout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := getSessionCookie(r)
	if err != nil {
		http.Error(w, "Invalid session: "+err.Error(), http.StatusUnauthorized)
		return
	}
	if err := p.store.EndSession(sessionID); err != nil {
		http.Error(w, "Cannot logout: "+err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, p.redirectURI, http.StatusSeeOther)
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

	idToken, err := p.verifier.Verify(p.ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	nonce, err := r.Cookie("nonce")
	if err != nil {
		http.Error(w, "Nonce not found", http.StatusBadRequest)
		return
	}
	if idToken.Nonce != nonce.Value {
		http.Error(w, "Nonce did not match", http.StatusBadRequest)
		return
	}

	var claims session.UserClaims
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	sessionID, err := p.store.NewSession(claimsToUser(&claims))
	if err != nil {
		http.Error(w, "Creating new session: "+err.Error(), http.StatusInternalServerError)
	}
	writeSessionCookieHeader(w, sessionID)
	http.Redirect(w, r, p.redirectURI, http.StatusSeeOther)
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

func claimsToUser(claims *session.UserClaims) *oapi.User {
	return &oapi.User{
		Sub:     claims.Sub,
		Iss:     claims.Iss,
		Name:    claims.Name,
		Email:   &claims.Email,
		Picture: &claims.Picture,
		Groups:  claims.Groups,
	}
}

func getUserContext(r *http.Request) (*oapi.User, error) {
	user, ok := r.Context().Value(contextKey).(*oapi.User)
	if !ok {
		return nil, fmt.Errorf("user context is incorrect type")
	}
	return user, nil
}
