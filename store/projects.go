package store

import (
	"fmt"

	"github.com/verifa/coastline/ent"
)

func (s *Store) QueryProjects() ([]*ent.Project, error) {
	projects, err := s.client.Project.Query().All(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("querying projects: %w", err)
	}

	return projects, nil
}

// func (s *Store) CreateProject() ([]*ent.Project, error) {
// 	s.client.Project.Create().SetName()
// 	projects, err := s.client.Project.
// 	if err != nil {
// 		return nil, fmt.Errorf("querying projects: %w", err)
// 	}

// 	return projects, nil
// }
