package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *PlanetStorage) CheckPlanetBelongsToUser(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (bool, error) {
	const getBuildingInfoQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM session_beta.planets p
			WHERE p.id = $1 AND p.user_id = $2
		);
	`

	var exists bool
	err := s.DB.QueryRow(ctx, getBuildingInfoQuery, planetID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("DB.QueryRow() %w", err)
	}

	return exists, nil
}
