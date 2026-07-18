package rating

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) getUsersRating(ctx context.Context, userID uuid.UUID) (models.Rating, error) {
	version, err := s.ratingStorage.GetLatestVersion(ctx, userRatingVersionType)
	if err != nil {
		return models.Rating{}, fmt.Errorf("GetLatestVersion(): %w", err)
	}

	topPlayers, err := s.ratingStorage.GetTopUsersRatingPlayers(ctx, version, topPlayersLimit)
	if err != nil {
		return models.Rating{}, fmt.Errorf("GetTopUsersRatingPlayers(): %w", err)
	}

	userRank, err := s.ratingStorage.GetUsersRatingRank(ctx, userID, version)
	if err != nil {
		if errors.Is(err, models.ErrUserNotInRating) {
			return models.Rating{
				Top: topPlayers,
			}, nil
		}

		return models.Rating{}, fmt.Errorf("GetUsersRatingRank(): %w", err)
	}

	var offset uint32
	if userRank > nearRankOffset {
		offset = userRank - nearRankOffset
	}

	nearPlayers, err := s.ratingStorage.GetUsersRatingPlayersByRankOffset(ctx, version, offset, nearPlayersLimit)
	if err != nil {
		return models.Rating{}, fmt.Errorf("GetUsersRatingPlayersByRankOffset(): %w", err)
	}

	return models.Rating{
		Top:  topPlayers,
		Near: nearPlayers,
	}, nil
}
