package planet

import (
	"context"
	"fmt"

	"initialservice/internal/models"

	"github.com/google/uuid"
)

// SetBuild sets the building level and stats for a given planet
func (s *PlanetStorage) SaveBuild(ctx context.Context, planetID uuid.UUID, buildInfo models.PlanetBuildInfo) error {
	const setBuildLvlQuery = `
		INSERT INTO session_beta.planet_buildings(
			planet_id,
			building_id, 
			updated_at,
			created_at
		) VALUES (
			$1,
			(SELECT id FROM session_beta.buildings WHERE type = $2 AND level = $3),
			$4,
			NOW()
		)
		ON CONFLICT (planet_id, building_id) DO UPDATE SET
			building_id = EXCLUDED.building_id,
			updated_at = EXCLUDED.updated_at;
		`

	_, err := s.DB.Exec(ctx, setBuildLvlQuery,
		planetID,
		buildInfo.Type,
		buildInfo.Level,
		buildInfo.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
