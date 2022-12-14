package store

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/request"
	"github.com/verifa/coastline/server/oapi"
	"github.com/verifa/coastline/worker"
)

func (s *Store) CreateTrigger(requestID uuid.UUID) (*oapi.Trigger, error) {
	dbTrigger, err := s.client.Trigger.Create().SetRequestID(requestID).Save(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("creating trigger in database: %w", err)
	}

	dbRequest, err := s.client.Request.Query().
		Where(request.ID(requestID)).
		WithProject().
		WithService().
		WithReviews().
		First(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("getting request with ID: %s: %w", requestID.String(), err)
	}

	msg := worker.TriggerMsg{
		TriggerID: dbTrigger.ID,
		Request:   dbRequestToAPI(dbRequest),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshalling request trigger message: %w", err)
	}

	// Publish a new trigger
	pubErr := s.nc.Publish(subjectTriggerRun, msgBytes)
	if pubErr != nil {
		return nil, fmt.Errorf("publishing request trigger: %w", pubErr)
	}
	return dbTriggerToAPI(dbTrigger), nil
}

func dbTriggerToAPI(dbTrigger *ent.Trigger) *oapi.Trigger {
	trigger := oapi.Trigger{
		Id: dbTrigger.ID,
	}
	workflows := make([]oapi.Workflow, len(dbTrigger.Edges.Workflows))

	for i, dbWorkflow := range dbTrigger.Edges.Workflows {
		workflows[i] = oapi.Workflow{
			Id:     dbWorkflow.ID,
			Output: dbWorkflow.Output,
			Error:  dbWorkflow.Error,
		}
	}
	trigger.Workflows = workflows

	return &trigger
}
