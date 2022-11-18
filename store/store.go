package store

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/hook"

	_ "github.com/xiaoqidun/entps"
)

type Config struct {
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
