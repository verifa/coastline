package server

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/verifa/coastline/server/oapi"
)

// populateDemoData populates the server with demo data
func (s *ServerImpl) populateDemoData() error {
	_, err := s.store.NewSession(devUser)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}
	project, err := s.store.CreateProject(&oapi.NewProject{
		Name: "silly-things",
	})
	if err != nil {
		return fmt.Errorf("creating project: %w", err)
	}
	for _, data := range demoDummyData {
		service, err := s.store.CreateService(data.service)
		if err != nil {
			return fmt.Errorf("creating service %s: %w", data.service.Name, err)
		}

		for _, req := range data.requests {
			req.ProjectId = project.Id
			req.ServiceId = service.Id
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(req)
			if err != nil {
				return fmt.Errorf("encoding request: %w", err)
			}
			vReq, err := s.engine.ValidateRequest(&buf)
			if err != nil {
				return fmt.Errorf("validating dummy request %s for service %s: %w", req.Kind, service.Name, err)
			}
			{
				_, err := s.store.CreateRequest(devUser, vReq)
				if err != nil {
					return fmt.Errorf("creating dummy request %s for service %s: %w", req.Kind, service.Name, err)
				}
			}

		}
	}

	return nil
}

type dummyData struct {
	service  *oapi.NewService
	requests []*oapi.NewRequest
}

var demoDummyData []*dummyData = []*dummyData{
	{
		service: &oapi.NewService{
			Name: "cat-facts",
			Labels: &oapi.NewService_Labels{
				AdditionalProperties: map[string]string{
					"tool": "cat-facts",
				},
			},
		},
		requests: []*oapi.NewRequest{
			{
				Kind: "CatFact",
				Spec: map[string]interface{}{
					"maxLength": 50,
				},
			},
		},
	},
	{
		service: &oapi.NewService{
			Name: "pokemon-info",
			Labels: &oapi.NewService_Labels{
				AdditionalProperties: map[string]string{
					"tool": "pokemon",
				},
			},
		},
		requests: []*oapi.NewRequest{
			{
				Kind: "PokemonInfo",
				Spec: map[string]interface{}{
					"name": "ditto",
				},
			},
			{
				Kind: "PokemonInfo",
				Spec: map[string]interface{}{
					"name": "eevee",
				},
			},
			{
				Kind: "PokemonInfo",
				Spec: map[string]interface{}{
					"name": "tyranitar",
				},
			},
		},
	},
}
