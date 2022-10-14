package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var projects = []project{
	{
		Name: "project-1",
	},
	{
		Name: "project-2",
	},
}

type project struct {
	Name string `json:"name"`
}

type projectsServer struct{}

func (p projectsServer) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", p.List)    // GET /projects - read a list of projects
	r.Post("/", p.Create) // POST /projects - create a new project and pepist it
	r.Put("/", p.Delete)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", p.Get)       // GET /projects/{id} - read a single project by :id
		r.Put("/", p.Update)    // PUT /projects/{id} - update a single project by :id
		r.Delete("/", p.Delete) // DELETE /projects/{id} - delete a single project by :id
		r.Get("/sync", p.Sync)
	})

	return r
}

func (p projectsServer) List(w http.ResponseWriter, r *http.Request) {
	// session, ok := r.Context().Value(contextKey).(Session)
	// if !ok {
	// 	http.Error(w, "Cannot get user from session context...", http.StatusInternalServerError)
	// 	return
	// }

	data := struct {
		Projects []project `json:"projects"`
	}{
		Projects: projects,
	}

	projectsJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Creating JSON from projects", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(projectsJSON)
}

func (p projectsServer) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("projects create"))
}

func (p projectsServer) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("project get"))
}

func (p projectsServer) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("project update"))
}

func (p projectsServer) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("project delete"))
}

func (p projectsServer) Sync(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("project sync"))
}
