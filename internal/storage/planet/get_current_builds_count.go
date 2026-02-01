package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *PlanetStorage) GetCurrentBuildsCount(ctx context.Context, planetID uuid.UUID) (uint8, error) {
	const getResourcesQuery = `
		SELECT count(id)
		FROM session_beta.building_events
		WHERE planet_id = $1;
	`

	var currentBuildsCount uint8
	err := r.DB.QueryRow(ctx, getResourcesQuery, planetID).Scan(
		&currentBuildsCount,
	)
	if err != nil {
		return 0, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return currentBuildsCount, nil
}
