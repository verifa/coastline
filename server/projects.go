package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/verifa/coastline/ent/project"
	"github.com/verifa/coastline/server/oapi"
)

func (s *ServerImpl) GetProjects(w http.ResponseWriter, r *http.Request, params oapi.GetProjectsParams) {
	resp, err := s.store.QueryProjects()
	if err != nil {
		http.Error(w, "Querying projects: "+err.Error(), http.StatusInternalServerError)
	}
	returnJSON(w, resp)
}

func (s *ServerImpl) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req oapi.NewProject
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	project, err := s.store.CreateProject(&req)
	if err != nil {
		http.Error(w, "Creating project: "+err.Error(), http.StatusBadRequest)
		return
	}
	returnJSON(w, project)
}

func (s *ServerImpl) GetProjectByID(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	resp, err := s.store.QueryProjects(project.ID(uuid.UUID(id)))
	if err != nil {
		http.Error(w, "Quering projects: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(resp.Projects) == 0 {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	returnJSON(w, resp.Projects[0])
}
