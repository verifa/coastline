package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/verifa/coastline/server/oapi"
	"github.com/verifa/coastline/store"
	"github.com/verifa/coastline/ui"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Config struct {
	DevMode        bool
	RedirectURI    string
	RequestsEngine RequestsEngineConfig
}

func DefaultConfig() Config {
	return Config{
		RequestsEngine: DefaultRequestsEngineConfig(),
		RedirectURI:    defaultEnv("CL_SERVER_REDIRECT_URI", "/ui"),
	}
}

func New(ctx context.Context, store *store.Store, config *Config) (*chi.Mux, error) {

	if config == nil {
		return nil, fmt.Errorf("config is required")
	}
	engine, err := NewRequestsEngine(&config.RequestsEngine)
	if err != nil {
		return nil, fmt.Errorf("creating requests engine: %w", err)
	}

	authProvider, err := newAuthProvider(ctx, config.DevMode, config.RedirectURI)
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

	serverImpl := ServerImpl{
		auth:   authProvider,
		store:  store,
		engine: engine,
	}
	wrapper := oapi.ServerInterfaceWrapper{
		Handler: &serverImpl,
	}

	// The handler produced by oapi-codegen is not very helpful when wanting to
	// implement authentication and authorization by middleware.
	// Instead of using all of oapi-codegen and hacking around it, it is easier
	// to use the parts of oapi-codegen which significantly help reduce boilerplate.
	// Therefore we mount the routes manually so that we have greater control of
	// which paths are protected via middleware and which are not.
	r.Route("/api/v1", func(r chi.Router) {

		// TODO: replace with OpenAPI generated wrapper by adding to spec
		r.Mount("/", authProvider.Routes())

		r.Group(func(r chi.Router) {
			// r.Use(authProvider.authenticateMiddleware)
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
			r.Get("/requestsspec", wrapper.GetRequestsSpec)
			//
			// UserInfo
			//
			r.Get("/userinfo", wrapper.GetUserInfo)
		})

		if oapiEnabled {
			r.Handle("/spec", http.StripPrefix("/api/v1/spec", http.FileServer(oapiSite)))
		}
	})

	// Setup frontend, if enabled (toggled via build tags)
	if ui.Enabled {
		// By default redirect the root path to ui
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/ui", http.StatusFound)
		})
		r.Mount("/ui", handleUI())
	}

	return r, nil
}

// handleUI returns a handler for our Single Page Application that checks if a
// requested resource exists, and if it doesn't, returns the root index.html
// (the single page).
func handleUI() http.Handler {
	index, err := ui.Site.Open("index.html")
	if err != nil {
		log.Fatal("Failed opening UI's index.html: " + err.Error())
	}
	var spaIndex bytes.Buffer
	if _, err := spaIndex.ReadFrom(index); err != nil {
		log.Fatal("Failed reading UI's index.html: " + err.Error())
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip the /ui prefix from the requested path to get the path to the
		// requested resource as it would be on the backend filesystem.
		path := strings.TrimPrefix(r.URL.Path, "/ui")
		// If requesting the root page, we will end up with nothing left, so
		// in that case we know it's the root page they were looking for
		if path == "" {
			w.WriteHeader(http.StatusAccepted)
			w.Write(spaIndex.Bytes())
			return
		}
		// Check if requested resource exists. If it does, treat it like a resource
		// such as a .js or .css file with the full path including the filename.
		// If it doesn't exist, it's a path without a filename and we should
		// return our Single Page (index.html)
		f, err := ui.Site.Open(path)
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusAccepted)
			w.Write(spaIndex.Bytes())
			return
		} else if err != nil {
			http.Error(w, "Error: opening requested path "+path+": "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		http.StripPrefix("/ui", http.FileServer(ui.Site)).ServeHTTP(w, r)
	})
}

var _ oapi.ServerInterface = (*ServerImpl)(nil)

type ServerImpl struct {
	auth   *authProvider
	store  *store.Store
	engine *RequestsEngine
}

func returnJSON(w http.ResponseWriter, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "Creating JSON response: "+err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(b)
}

func returnBytesAsJSON(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(b)
}

func defaultEnv(env string, value string) string {
	e, ok := os.LookupEnv(env)
	if ok {
		return e
	}
	return value
}
