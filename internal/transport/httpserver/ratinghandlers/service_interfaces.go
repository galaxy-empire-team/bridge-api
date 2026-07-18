package ratinghandlers

import (
	"context"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/google/uuid"
)

type RatingService interface {
	GetRating(ctx context.Context, userID uuid.UUID) (models.Ratings, error)
}
