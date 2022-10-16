package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/verifa/coastline/ent/request"
	"github.com/verifa/coastline/server/oapi"
)

func (s *ServerImpl) GetRequests(w http.ResponseWriter, r *http.Request, params oapi.GetRequestsParams) {
	resp, err := s.store.QueryRequests()
	if err != nil {
		http.Error(w, "Querying requests: "+err.Error(), http.StatusInternalServerError)
	}
	returnJSON(w, resp)
}

func (s *ServerImpl) CreateRequest(w http.ResponseWriter, r *http.Request) {
	var req oapi.NewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	request, err := s.store.CreateRequest(&req)
	if err != nil {
		http.Error(w, "Creating request: "+err.Error(), http.StatusBadRequest)
		return
	}
	returnJSON(w, request)
}

func (s *ServerImpl) GetRequestByID(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	resp, err := s.store.QueryRequests(request.ID(uuid.UUID(id)))
	if err != nil {
		http.Error(w, "Quering requests: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(resp.Requests) == 0 {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}
	returnJSON(w, resp.Requests[0])
}
