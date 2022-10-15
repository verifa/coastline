package store

import (
	"context"
	"fmt"

	"github.com/verifa/coastline/ent"

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
	if err := s.init(); err != nil {
		return nil, fmt.Errorf("initializing database: %w", err)
	}

	return &s, nil
}

type Store struct {
	client *ent.Client
	ctx    context.Context
}

func (s *Store) Close() error {
	return s.client.Close()
}

// init populates the database with some initial data
func (s *Store) init() error {
	_, err := s.client.Project.Create().SetName("default").Save(s.ctx)
	if err != nil {
		return fmt.Errorf("creating project: %w", err)
	}

	return nil
}
