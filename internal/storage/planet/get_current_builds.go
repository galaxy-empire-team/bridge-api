package planet

import (
	"context"
	"fmt"
	"time"

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
	var startedAt, finishedAt time.Time
	for rows.Next() {
		var b models.BuildingInProgress
		if err := rows.Scan(&b.BuildingID, &startedAt, &finishedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		b.StartedAt = startedAt.UTC()
		b.FinishedAt = finishedAt.UTC()

		builds = append(builds, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return builds, nil
}
