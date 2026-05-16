package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

// CheckPlanetOwner checks if the planet is owned by the user.
// Returns ErrPlanetDoesNotBelongToUser if not.
func (r *Repository) CheckPlanetOwner(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error {
	if userID == uuid.Nil {
		return models.ErrNoUserIDProvided
	}

	if planetID == uuid.Nil {
		return models.ErrNoPlanetIDProvided
	}

	isUserPlanet, err := r.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetID)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}

	if !isUserPlanet {
		return models.ErrPlanetDoesNotBelongToUser
	}

	return nil
}
