package server

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/verifa/coastline/server/oapi"
)

const cookieName = "coastline_session"

func getSessionCookie(r *http.Request) (uuid.UUID, error) {
	sessionCookie, err := r.Cookie(cookieName)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.Parse(sessionCookie.Value)
}

func writeSessionCookieHeader(w http.ResponseWriter, sessionID uuid.UUID) {
	// Create cookie
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    sessionID.String(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Hour.Seconds() * 8),
		Path:     "/",
	})
}

func (s *ServerImpl) GetUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: move this to store package and make a dbUserToAPI function to convert
	// ent.User to oapi.User
	dbUsers, err := s.store.Client().User.Query().All(s.auth.ctx)
	if err != nil {
		http.Error(w, "Getting users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := oapi.UsersResp{
		Users: make([]oapi.User, len(dbUsers)),
	}
	for i, dbUser := range dbUsers {
		resp.Users[i] = oapi.User{
			Sub:     dbUser.Sub,
			Iss:     dbUser.Iss,
			Name:    dbUser.Name,
			Email:   dbUser.Email,
			Picture: dbUser.Picture,
		}
	}
	returnJSON(w, resp)
}
