package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetCurrentBuilds(ctx context.Context, planetID uuid.UUID) ([]models.BuildingInProgress, error) {
	const getCurrentBuildsQuery = `
		SELECT 
			building_id,
			started_at,
			finished_at
		FROM session_beta.event_buildings
		WHERE planet_id = $1;
	`

	rows, err := r.DB.Query(ctx, getCurrentBuildsQuery, planetID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query(): %w", err)
	}
	defer rows.Close()

	var builds []models.BuildingInProgress
	for rows.Next() {
		var b models.BuildingInProgress
		if err := rows.Scan(&b.BuildingID, &b.StartedAt, &b.FinishedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		builds = append(builds, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return builds, nil
}
