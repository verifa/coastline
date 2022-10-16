package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/verifa/coastline/server/oapi"
	"github.com/verifa/coastline/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func New(ctx context.Context, store *store.Store) (*chi.Mux, error) {

	provider, err := newAuthProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating authentication provider: %w", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/healthz"))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"*"},
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:9998"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		Debug:            true,
	}))

	// r.Get("/api/v1/swagger.json", func(w http.ResponseWriter, r *http.Request) {
	// 	swagger, err := oapi.GetSwagger()
	// 	if err != nil {
	// 		http.Error(w, "Generating swagger: "+err.Error(), http.StatusInternalServerError)
	// 	}
	// 	returnJSON(w, swagger)
	// })

	// TODO: OpenAPI docs
	// r.Handle("/swaggerui/*",
	// 	http.StripPrefix("/swaggerui", http.FileServer(http.FS(swaggerDist))),
	// )
	// handler := oapi.HandlerWithOptions(&ServerImpl{store: store}, oapi.ChiServerOptions{
	// 	Middlewares: []oapi.MiddlewareFunc{
	// 		// func(hf http.HandlerFunc) http.HandlerFunc {
	// 		// 	return func(w http.ResponseWriter, r *http.Request) {
	// 		// 		http.Error(w, "no session", http.StatusUnauthorized)
	// 		// 	}
	// 		// },
	// 	},
	// })

	wrapper := oapi.ServerInterfaceWrapper{
		Handler: &ServerImpl{store: store},
	}

	// The handler produced by oapi-codegen is not very helpful when wanting to
	// implement authentication and authorization by middleware.
	// Instead of using all of oapi-codegen and hacking around it, it is easier
	// to use the parts of oapi-codegen which significantly help reduce boilerplate.
	// Therefore we mount the routes manually so that we have greater control of
	// which paths are protected via middleware and which are not.
	r.Route("/api/v1", func(r chi.Router) {

		// TODO: replace with OpenAPI generated wrapper by adding to spec
		r.Mount("/", provider.Routes())
		r.Group(func(r chi.Router) {
			// r.Use(provider.authenticateMiddleware)
			//
			// Projects
			//
			r.Get("/projects", wrapper.GetProjects)
			r.Post("/projects", wrapper.CreateProject)
			r.Post("/projects/{id}", wrapper.GetProjectByID)
			//
			// Services
			//
			r.Get("/services", wrapper.GetServices)
			r.Post("/services", wrapper.CreateService)
			r.Post("/services/{id}", wrapper.GetServiceByID)
			//
			// Requests
			//
			r.Get("/requests", wrapper.GetRequests)
			r.Post("/requests", wrapper.CreateRequest)
			r.Post("/requests/{id}", wrapper.GetRequestByID)
		})
	})

	return r, nil
}

var _ oapi.ServerInterface = (*ServerImpl)(nil)

type ServerImpl struct {
	store *store.Store
}

func returnJSON(w http.ResponseWriter, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "Creating JSON response: "+err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(b)
}
