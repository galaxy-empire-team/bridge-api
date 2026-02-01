package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// GetBuildLvl retrieves mine infromation from the target planet's buildings.
func (s *PlanetStorage) GetBuildingLvl(ctx context.Context, planetID uuid.UUID, BuildingType consts.BuildingType) (uint8, error) {
	const getBuildingInfoQuery = `
		SELECT b.level
		FROM session_beta.planet_buildings pb
		JOIN session_beta.buildings b ON pb.building_id = b.id
		WHERE pb.planet_id = $1 AND b.building_type = $2;
	`

	var buildLvl uint8
	err := s.DB.QueryRow(ctx, getBuildingInfoQuery, planetID, BuildingType).Scan(&buildLvl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return buildLvl, nil
}
