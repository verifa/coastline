package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/verifa/coastline/server/oapi"
)

func (s *ServerImpl) ReviewRequest(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var req oapi.NewReview
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := getUserContext(r)
	if err != nil {
		http.Error(w, "Getting user context: "+err.Error(), http.StatusInternalServerError)
	}

	review, err := s.store.CreateReview(id, user, &req)
	if err != nil {
		http.Error(w, "Creating review: "+err.Error(), http.StatusBadRequest)
		return
	}
	returnJSON(w, review)
}
