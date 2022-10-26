package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/verifa/coastline/ent/request"
	"github.com/verifa/coastline/server/oapi"
)

func (s *ServerImpl) GetRequests(w http.ResponseWriter, r *http.Request, params oapi.GetRequestsParams) {
	resp, err := s.store.QueryRequests()
	if err != nil {
		http.Error(w, "Querying requests: "+err.Error(), http.StatusInternalServerError)
		return
	}
	returnJSON(w, resp)
}

func (s *ServerImpl) CreateRequest(w http.ResponseWriter, r *http.Request) {
	var req oapi.NewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Validate the request against the requests engine
	if err := s.engine.Validate(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
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

func (s *ServerImpl) GetRequestTemplateSpec(w http.ResponseWriter, r *http.Request, id string) {
	spec, err := s.engine.OpenAPISpec(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Generating OpenAPI specification for %s: %s", id, err.Error()), http.StatusBadRequest)
	}
	returnBytesAsJSON(w, spec)
}
