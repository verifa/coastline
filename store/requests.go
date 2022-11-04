package store

import (
	"encoding/json"
	"fmt"

	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/predicate"
	"github.com/verifa/coastline/ent/request"
	"github.com/verifa/coastline/ent/review"
	"github.com/verifa/coastline/server/oapi"
	"github.com/verifa/coastline/worker"
)

func (s *Store) QueryRequests(ps ...predicate.Request) (*oapi.RequestsResp, error) {
	dbRequests, err := s.client.Request.Query().Where(ps...).
		WithProject().
		WithService().
		WithReviews().
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

	dbRequest, err := s.client.Request.Create().
		SetKind(req.Kind).
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

	c := m.Client()

	dbTrigger, err := c.Trigger.Create().SetRequestID(requestID).Save(s.ctx)
	if err != nil {
		return fmt.Errorf("creating trigger in database: %w", err)
	}

	dbRequest, err := c.Request.Query().
		Where(request.ID(requestID)).
		WithProject().
		WithService().
		WithReviews().
		First(s.ctx)
	if err != nil {
		return fmt.Errorf("getting request with ID: %s: %w", requestID.String(), err)
	}

	msg := worker.TriggerMsg{
		TriggerID: dbTrigger.ID,
		Request:   dbRequestToAPI(dbRequest),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshalling request trigger message: %w", err)
	}

	// Publish a new trigger
	pubErr := s.nc.Publish(subjectTriggerRun, msgBytes)
	if pubErr != nil {
		return fmt.Errorf("publishing request trigger: %w", pubErr)
	}

	return nil
}

func dbRequestToAPI(dbRequest *ent.Request) *oapi.Request {
	request := oapi.Request{
		Id:     dbRequest.ID,
		Kind:   dbRequest.Kind,
		Status: oapi.RequestStatus(dbRequest.Status),
		Spec:   dbRequest.Spec,
	}
	if dbRequest.Edges.Project != nil {
		request.Project = *dbProjectToAPI(dbRequest.Edges.Project)
	}
	if dbRequest.Edges.Service != nil {
		request.Service = *dbServiceToAPI(dbRequest.Edges.Service)
	}
	if dbRequest.Edges.Reviews != nil {
		for _, dbReview := range dbRequest.Edges.Reviews {
			request.Reviews = append(request.Reviews, *dbReviewToAPI(dbReview))
		}
	}
	return &request
}
