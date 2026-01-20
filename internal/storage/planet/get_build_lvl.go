package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"initialservice/internal/models"
)

// GetBuildLvl retrieves mine infromation from the target planet's buildings
func (s *PlanetStorage) GetBuildLvl(ctx context.Context, planetID uuid.UUID, buildType models.BuildType) (uint8, error) {
	const getMineInfoQuery = `
		SELECT b.level
		FROM session_beta.planet_buildings pb
		JOIN session_beta.buildings b ON pb.building_id = b.id
		WHERE pb.planet_id = $1 AND b.type = $2;
	`

	var buildLvl uint8
	err := s.DB.QueryRow(ctx, getMineInfoQuery, planetID, buildType).Scan(&buildLvl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return buildLvl, nil
}
