package store

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/hook"
	"github.com/verifa/coastline/server/oapi"

	_ "github.com/xiaoqidun/entps"
)

type Config struct {
	// InitData specifies whether to load the database with dummy data or not.
	// It is intended for demoing purposes and should be ignored in production
	InitData bool

	// SkipHook registers store hooks if enabled. This is for testing purposes.
	SkipHooks bool
	// SkipNATS connects to NATS if enabled. This is for testing purposes.
	SkipNATS bool

	SessionDuration time.Duration
}

func New(ctx context.Context, config *Config) (*Store, error) {
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if config.SessionDuration == 0 {
		// Set default duration
		config.SessionDuration = time.Hour * 2
	}

	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("creating schema: %w", err)
	}

	var nc *nats.Conn
	if !config.SkipNATS {
		var err error
		nc, err = setupNATS()
		if err != nil {
			return nil, fmt.Errorf("setting up nats: %w", err)
		}
	}

	s := Store{
		ctx:    ctx,
		config: config,
		client: client,
		nc:     nc,
	}

	if !config.SkipNATS {
		// Subscribe to NATS
		s.natsSubscribe()
	}

	if !config.SkipHooks {
		// Setup hooks
		s.RegisterHooks()
	}

	if config.InitData {
		if err := s.init(); err != nil {
			return nil, fmt.Errorf("initializing database: %w", err)
		}
	}

	return &s, nil
}

type Store struct {
	ctx    context.Context
	config *Config
	client *ent.Client
	nc     *nats.Conn
}

func (s *Store) Client() *ent.Client {
	return s.client
}

func (s *Store) Close() error {
	return s.client.Close()
}

func (s *Store) RegisterHooks() {
	s.client.Review.Use(func(next ent.Mutator) ent.Mutator {
		return hook.ReviewFunc(func(ctx context.Context, m *ent.ReviewMutation) (ent.Value, error) {
			// Execute mutation first
			value, err := next.Mutate(ctx, m)
			if err != nil {
				return nil, err
			}
			if err := s.HandleNewReview(m); err != nil {
				return nil, fmt.Errorf("handling new review hook: %w", err)
			}
			return value, nil
		})
	})
	s.client.Request.Use(func(next ent.Mutator) ent.Mutator {
		return hook.RequestFunc(func(ctx context.Context, m *ent.RequestMutation) (ent.Value, error) {
			// Execute mutation first
			value, err := next.Mutate(ctx, m)
			if err != nil {
				return nil, err
			}
			if err := s.HandleUpdatedRequest(m); err != nil {
				return nil, fmt.Errorf("handling updated request hook: %w", err)
			}
			return value, nil
		})
	})
}

// init populates the database with some initial data
func (s *Store) init() error {

	dummyUser := oapi.User{
		Sub:  "dummy",
		Iss:  "dummy",
		Name: "dummy",
	}
	// TODO: this is really hacky as it doesn't check if the data already exists
	// so if using persistent data it will fail...
	_, err := s.createUpdateUser(&dummyUser)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}
	project, err := s.CreateProject(&oapi.NewProject{
		Name: "dummy-project",
	})
	if err != nil {
		return fmt.Errorf("creating project: %w", err)
	}
	service, err := s.CreateService(&oapi.NewService{
		Name: "dummy-service",
		Labels: &oapi.NewService_Labels{
			AdditionalProperties: map[string]string{
				"tool": "artifactory",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("creating service: %w", err)
	}
	{
		requests := []*oapi.NewRequest{
			{
				ProjectId: project.Id,
				ServiceId: service.Id,
				Kind:      "JenkinsServerRequest",
				Spec: map[string]interface{}{
					"name": "server-1",
				},
			},
			{
				ProjectId: project.Id,
				ServiceId: service.Id,
				Kind:      "JenkinsServerRequest",
				Spec: map[string]interface{}{
					"name": "server-2",
				},
			},
			{
				ProjectId: project.Id,
				ServiceId: service.Id,
				Kind:      "JenkinsServerRequest",
				Spec: map[string]interface{}{
					"name": "server-3",
				},
			},
		}
		for i, req := range requests {
			if _, err := s.CreateRequest(&dummyUser, req); err != nil {
				return fmt.Errorf("creating request %d: %w", i, err)
			}
		}
	}

	return nil
}
