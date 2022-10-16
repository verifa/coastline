package store

import (
	"fmt"

	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/predicate"
	"github.com/verifa/coastline/server/oapi"
)

func (s *Store) QueryRequests(ps ...predicate.Request) (*oapi.RequestsResp, error) {
	dbRequests, err := s.client.Request.Query().Where(ps...).All(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("querying requests: %w", err)
	}

	var requests []oapi.Request
	for _, dbRequest := range dbRequests {
		requests = append(requests, dbRequestToAPI(dbRequest))
	}

	return &oapi.RequestsResp{
		Requests: requests,
	}, nil
}

func (s *Store) CreateRequest(req *oapi.NewRequest) (*oapi.Request, error) {
	dbRequest, err := s.client.Request.Create().
		SetType(req.Type).
		SetProjectID(req.ProjectId).
		SetServiceID(req.ServiceId).
		SetRequestedBy(req.RequestedBy).
		SetSpec(req.Spec).
		Save(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	request := dbRequestToAPI(dbRequest)
	return &request, nil
}

func dbRequestToAPI(dbRequest *ent.Request) oapi.Request {
	return oapi.Request{
		Id:          dbRequest.ID,
		Type:        dbRequest.Type,
		RequestedBy: dbRequest.RequestedBy,
		Spec:        dbRequest.Spec,
		// TODO: handle project_id and service_id
	}
}
