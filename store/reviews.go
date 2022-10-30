package store

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/review"
	"github.com/verifa/coastline/server/oapi"
)

func (s *Store) CreateReview(requestID uuid.UUID, user *oapi.User, req *oapi.NewReview) (*oapi.Review, error) {
	dbUser, err := s.getEntUser(user)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}

	dbReview, err := s.client.Review.Create().
		SetRequestID(requestID).
		SetStatus(review.Status(req.Status)).
		SetType(review.Type(req.Type)).
		SetCreatedBy(dbUser).
		Save(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("saving review: %w", err)
	}
	return dbReviewToAPI(dbReview), nil
}

func dbReviewToAPI(dbReview *ent.Review) *oapi.Review {
	review := oapi.Review{
		Id:     dbReview.ID,
		Status: oapi.ReviewStatus(dbReview.Status),
		Type:   oapi.ReviewType(dbReview.Type),
	}
	if dbReview.Edges.CreatedBy != nil {
		review.CreatedBy = *dbUserToAPI(dbReview.Edges.CreatedBy)
	}

	return &review
}
