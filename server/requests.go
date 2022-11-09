package server

import (
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

	req, err := s.engine.ValidateRequest(r.Body)
	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := getUserContext(r)
	if err != nil {
		http.Error(w, "Getting user context: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request, err := s.store.CreateRequest(user, req)
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
