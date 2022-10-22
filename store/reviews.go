package store

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/review"
	"github.com/verifa/coastline/server/oapi"
)

func (s *Store) CreateReview(requestID uuid.UUID, req *oapi.NewReview) (*oapi.Review, error) {
	dbReview, err := s.client.Review.Create().
		SetRequestID(requestID).
		SetStatus(review.Status(req.Status)).
		SetType(review.Type(req.Type)).
		Save(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("saving review: %w", err)
	}
	return &oapi.Review{
		Id:     dbReview.ID,
		Status: oapi.ReviewStatus(dbReview.Status),
		Type:   oapi.ReviewType(dbReview.Type),
	}, nil
}

func dbReviewToAPI(dbReview *ent.Review) *oapi.Review {
	return &oapi.Review{
		Id:     dbReview.ID,
		Status: oapi.ReviewStatus(dbReview.Status),
		Type:   oapi.ReviewType(dbReview.Type),
	}
}
