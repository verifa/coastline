package store

import (
	"fmt"

	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/predicate"
	"github.com/verifa/coastline/server/oapi"
)

func (s *Store) QueryServices(ps ...predicate.Service) (*oapi.ServicesResp, error) {
	dbServices, err := s.client.Service.Query().Where(ps...).All(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("querying services: %w", err)
	}

	var services = make([]oapi.Service, len(dbServices))
	for i, dbService := range dbServices {
		services[i] = dbServiceToAPI(dbService)
	}

	return &oapi.ServicesResp{
		Services: services,
	}, nil
}

func (s *Store) CreateService(req *oapi.NewService) (*oapi.Service, error) {
	dbService, err := s.client.Service.Create().
		SetName(req.Name).
		Save(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("creating service: %w", err)
	}
	service := dbServiceToAPI(dbService)
	return &service, nil
}

func dbServiceToAPI(dbService *ent.Service) oapi.Service {
	return oapi.Service{
		Id:   dbService.ID,
		Name: dbService.Name,
	}
}
