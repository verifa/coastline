package store

import (
	"context"
	"fmt"

	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/hook"
	"github.com/verifa/coastline/server/oapi"

	_ "github.com/xiaoqidun/entps"
)

func New(ctx context.Context) (*Store, error) {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("creating schema: %w", err)
	}

	s := Store{
		client: client,
		ctx:    ctx,
	}

	// Setup hooks
	s.RegisterHooks()

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("initializing database: %w", err)
	}

	return &s, nil
}

type Store struct {
	client *ent.Client
	ctx    context.Context
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
}

// init populates the database with some initial data
func (s *Store) init() error {
	// TODO: this is really hacky as it doesn't check if the data already exists
	// so if using persistent data it will fail...
	project, err := s.CreateProject(&oapi.NewProject{
		Name: "dummy-project",
	})
	if err != nil {
		return fmt.Errorf("creating project: %w", err)
	}
	service, err := s.CreateService(&oapi.NewService{
		Name: "dummy-service",
	})
	if err != nil {
		return fmt.Errorf("creating service: %w", err)
	}
	{
		requests := []*oapi.NewRequest{
			{
				ProjectId:   project.Id,
				ServiceId:   service.Id,
				RequestedBy: "dummy-user",
				Type:        "JenkinsServerRequest",
				Spec: map[string]interface{}{
					"name": "server-1",
				},
			},
			{
				ProjectId:   project.Id,
				ServiceId:   service.Id,
				RequestedBy: "dummy-user",
				Type:        "JenkinsServerRequest",
				Spec: map[string]interface{}{
					"name": "server-2",
				},
			},
			{
				ProjectId:   project.Id,
				ServiceId:   service.Id,
				RequestedBy: "dummy-user",
				Type:        "JenkinsServerRequest",
				Spec: map[string]interface{}{
					"name": "server-3",
				},
			},
		}
		for i, req := range requests {
			if _, err := s.CreateRequest(req); err != nil {
				return fmt.Errorf("creating request %d: %w", i, err)
			}
		}
	}

	return nil
}
