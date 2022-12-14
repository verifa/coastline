package store

import (
	"fmt"

	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/predicate"
	"github.com/verifa/coastline/ent/request"
	"github.com/verifa/coastline/ent/review"
	"github.com/verifa/coastline/server/oapi"
)

func (s *Store) QueryRequests(ps ...predicate.Request) (*oapi.RequestsResp, error) {
	dbRequests, err := s.client.Request.Query().Where(ps...).
		WithProject().
		WithService().
		WithReviews().
		WithTriggers(func(tq *ent.TriggerQuery) {
			tq.WithWorkflows()
		}).
		All(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("querying requests: %w", err)
	}

	var requests = make([]oapi.Request, len(dbRequests))
	for i, dbRequest := range dbRequests {
		requests[i] = *dbRequestToAPI(dbRequest)
	}

	return &oapi.RequestsResp{
		Requests: requests,
	}, nil
}

func (s *Store) CreateRequest(user *oapi.User, req *oapi.NewRequest) (*oapi.Request, error) {

	dbUser, err := s.getEntUser(user)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}

	var descr string
	if req.Description != nil {
		descr = *req.Description
	}

	dbRequest, err := s.client.Request.Create().
		SetKind(req.Kind).
		SetDescription(descr).
		SetProjectID(req.ProjectId).
		SetServiceID(req.ServiceId).
		SetSpec(req.Spec).
		SetCreatedBy(dbUser).
		Save(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	return dbRequestToAPI(dbRequest), nil
}

// HandleNewReview is called whenever a new review is added to a request and is
// responsible for re-evaluating the status of a request
func (s *Store) HandleNewReview(m *ent.ReviewMutation) error {
	requestID, ok := m.RequestID()
	if !ok {
		id, _ := m.ID()
		return fmt.Errorf("no request for review with ID: %s", id)
	}
	c := m.Client()

	// Current logic for this is really bad and needs re-work.
	// Right now it checks if there are any reviews that reject and sets status
	// to rejected. If any approve, and no rejects, set to approve
	dbReviews, err := c.Review.Query().
		Where(review.HasRequestWith(request.ID(requestID))).
		// Order by date time
		Order(ent.Desc(review.FieldCreateTime)).
		All(s.ctx)
	if err != nil {
		return fmt.Errorf("getting reviews: %w", err)
	}

	var approve bool
	for _, r := range dbReviews {
		if r.Status == review.StatusApprove {
			approve = true
		} else {
			approve = false
			break
		}
	}
	if len(dbReviews) > 0 {
		var status request.Status
		if approve {
			status = request.StatusApproved
		} else {
			status = request.StatusRejected
		}
		_, err := c.Request.UpdateOneID(requestID).
			SetStatus(status).
			Save(s.ctx)
		if err != nil {
			return fmt.Errorf("updating request with ID %s: %w", requestID, err)
		}
	}
	return nil
}

// HandleUpdatedRequest is called whenever a request is updated (or mutated, in ent terms)
// and is responsible for triggering a deployment if the request is approved
func (s *Store) HandleUpdatedRequest(m *ent.RequestMutation) error {
	requestStatus, ok := m.Status()
	if !ok {
		// Shouldn't happen, but let's just exit. Or maybe we should error?
		return nil
	}
	// For now we only care about the Approved status.
	// Finish here if it's not the Approved status
	if requestStatus != request.StatusApproved {
		return nil
	}

	requestID, ok := m.ID()
	if !ok {
		return fmt.Errorf("request does not have ID")
	}

	_, err := s.CreateTrigger(requestID)
	if err != nil {
		return fmt.Errorf("creating trigger: %w", err)
	}

	return nil
}

func dbRequestToAPI(dbRequest *ent.Request) *oapi.Request {
	request := oapi.Request{
		Id:          dbRequest.ID,
		Kind:        dbRequest.Kind,
		Description: dbRequest.Description,
		Status:      oapi.RequestStatus(dbRequest.Status),
		Spec:        dbRequest.Spec,
		Triggers:    make([]oapi.Trigger, len(dbRequest.Edges.Triggers)),
		Reviews:     make([]oapi.Review, len(dbRequest.Edges.Reviews)),
	}
	if dbRequest.Edges.Project != nil {
		request.Project = *dbProjectToAPI(dbRequest.Edges.Project)
	}
	if dbRequest.Edges.Service != nil {
		request.Service = *dbServiceToAPI(dbRequest.Edges.Service)
	}
	for i, dbReview := range dbRequest.Edges.Reviews {
		request.Reviews[i] = *dbReviewToAPI(dbReview)
	}
	for i, dbTrigger := range dbRequest.Edges.Triggers {
		request.Triggers[i] = *dbTriggerToAPI(dbTrigger)
	}
	return &request
}
