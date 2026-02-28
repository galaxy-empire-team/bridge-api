package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *PlanetStorage) GetBuildsInProgressCount(ctx context.Context, planetID uuid.UUID) (uint8, error) {
	const getCurrentBuildsCountQuery = `
		SELECT count(id)
		FROM session_beta.event_buildings
		WHERE planet_id = $1;
	`

	var currentBuildsCount uint8
	err := r.DB.QueryRow(ctx, getCurrentBuildsCountQuery, planetID).Scan(
		&currentBuildsCount,
	)
	if err != nil {
		return 0, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return currentBuildsCount, nil
}
