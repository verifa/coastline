package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/verifa/coastline/ent/service"
	"github.com/verifa/coastline/server/oapi"
)

func (s *ServerImpl) GetServices(w http.ResponseWriter, r *http.Request, params oapi.GetServicesParams) {
	resp, err := s.store.QueryServices()
	if err != nil {
		http.Error(w, "Querying services: "+err.Error(), http.StatusInternalServerError)
	}
	returnJSON(w, resp)
}

func (s *ServerImpl) CreateService(w http.ResponseWriter, r *http.Request) {
	var req oapi.NewService
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	service, err := s.store.CreateService(&req)
	if err != nil {
		http.Error(w, "Creating service: "+err.Error(), http.StatusBadRequest)
		return
	}
	returnJSON(w, service)
}

func (s *ServerImpl) GetServiceByID(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	resp, err := s.store.QueryServices(service.ID(uuid.UUID(id)))
	if err != nil {
		http.Error(w, "Quering services: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(resp.Services) == 0 {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}
	returnJSON(w, resp.Services[0])
}
