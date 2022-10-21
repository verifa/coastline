package server

import (
	"net/http"
)

func (s ServerImpl) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	ui, err := s.auth.UserInfo(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	returnJSON(w, ui)
}
