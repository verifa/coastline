package store

import (
	"fmt"

	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/predicate"
	"github.com/verifa/coastline/server/oapi"
)

func (s *Store) QueryProjects(ps ...predicate.Project) (*oapi.ProjectsResp, error) {
	dbProjects, err := s.client.Project.Query().Where(ps...).All(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("querying projects: %w", err)
	}

	var projects []oapi.Project
	for _, dbProject := range dbProjects {
		projects = append(projects, dbProjectToAPI(dbProject))
	}

	return &oapi.ProjectsResp{
		Projects: projects,
	}, nil
}

func (s *Store) CreateProject(newProject *oapi.NewProject) (*oapi.Project, error) {
	dbProject, err := s.client.Project.Create().
		SetName(newProject.Name).
		Save(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("creating project: %w", err)
	}
	project := dbProjectToAPI(dbProject)
	return &project, nil
}

func dbProjectToAPI(dbProject *ent.Project) oapi.Project {
	return oapi.Project{
		Id:   dbProject.ID,
		Name: dbProject.Name,
	}
}
