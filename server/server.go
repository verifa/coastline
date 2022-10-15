package server

import (
	"context"
	"fmt"

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

	r.Mount("/api/v1", provider.Routes())

	r.Group(func(r chi.Router) {
		// r.Use(provider.authorizeMiddleware)
		r.Mount("/api/v1/projects", projectsServer{store: store}.Routes())
	})

	return r, nil
}
