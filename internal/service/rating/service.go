package rating

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

const (
	userRatingVersionType  = "economy_points"
	fleetRatingVersionType = "fleet_points"
	topPlayersLimit        = 30
	nearRankOffset         = 3
	nearPlayersLimit       = 5
)

type ratingStorage interface {
	GetLatestVersion(ctx context.Context, versionType string) (uint32, error)
	GetTopUsersRatingPlayers(ctx context.Context, version uint32, limit uint32) ([]models.RatingPlayer, error)
	GetUsersRatingRank(ctx context.Context, userID uuid.UUID, version uint32) (uint32, error)
	GetUsersRatingPlayersByRankOffset(ctx context.Context, version uint32, offset uint32, limit uint32) ([]models.RatingPlayer, error)
	GetTopFleetRatingPlayers(ctx context.Context, version uint32, limit uint32) ([]models.RatingPlayer, error)
	GetFleetRatingRank(ctx context.Context, userID uuid.UUID, version uint32) (uint32, error)
	GetFleetRatingPlayersByRankOffset(ctx context.Context, version uint32, offset uint32, limit uint32) ([]models.RatingPlayer, error)
}

type Service struct {
	ratingStorage ratingStorage
}

func New(ratingStorage ratingStorage) *Service {
	return &Service{ratingStorage: ratingStorage}
}
