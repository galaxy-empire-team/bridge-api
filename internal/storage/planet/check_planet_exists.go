package planet

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *PlanetStorage) CheckPlanetExists(ctx context.Context, coordinates models.Coordinates) (bool, error) {
	const getBuildingInfoQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM session_beta.planets p
			WHERE x = $1 AND y = $2 AND z = $3
		);
	`

	var exists bool
	err := s.DB.QueryRow(ctx, getBuildingInfoQuery, coordinates.X, coordinates.Y, coordinates.Z).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return exists, nil
}
