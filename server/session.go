package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/google/uuid"
)

type Sessioner struct {
	cookieName string
	sessions   map[string]*Session
}

type Session struct {
	ID      uuid.UUID
	IDToken *oidc.IDToken

	UserInfo UserInfo
}

type UserInfo struct {
	Name   string   `json:"name"`
	Email  string   `json:"email"`
	Groups []string `json:"groups"`
}

type UserClaims struct {
	UserID string   `json:"sub"`
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
}

func newSessioner() *Sessioner {
	return &Sessioner{
		sessions:   make(map[string]*Session),
		cookieName: "project_session",
	}
}

func (s *Sessioner) NewSession(w http.ResponseWriter, claims *UserClaims) {

	sessionID := uuid.New()
	session := Session{
		ID: sessionID,
		// IDToken: idToken,
	}

	session.UserInfo = UserInfo{
		Name:   claims.Name,
		Email:  claims.Email,
		Groups: claims.Groups,
	}
	s.sessions[sessionID.String()] = &session
	http.SetCookie(w, &http.Cookie{
		Name:     s.cookieName,
		Value:    sessionID.String(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Hour.Seconds()),
		Path:     "/",
	})
}

func (s *Sessioner) EndSession(r *http.Request) error {
	session, err := s.AuthorizeSession(r)
	if err != nil {
		return err
	}
	delete(s.sessions, session.ID.String())

	return nil
}

func (s *Sessioner) AuthorizeSession(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie(s.cookieName)
	if err != nil {
		return nil, err
	}
	session, ok := s.sessions[sessionCookie.Value]
	if !ok {
		return nil, errors.New("no session exists")
	}
	return session, nil
}

// func (s *Sessioner) Get(id string) (Session, bool) {
// 	session, ok := s.sessions[id]
// 	return session, ok
// }

// func newSession() (*session.Manager, error) {
// 	mgr, err := session.NewManager("memory", session.Options{
// 		CookieName: cookieName,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("creating session manager: %w", err)
// 	}
// 	return mgr, nil
// }
