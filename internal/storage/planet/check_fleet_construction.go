package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *PlanetStorage) CheckFleetConstruction(ctx context.Context, planetID uuid.UUID) (bool, error) {
	const getFleetConstructionQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM session_beta.event_fleet_constructions p
			WHERE planet_id = $1
		);
	`

	var exists bool
	err := s.DB.QueryRow(ctx, getFleetConstructionQuery, planetID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return exists, nil
}
