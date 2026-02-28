package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// GetBuildingID retrieves mine infromation from the target planet's building.
func (s *PlanetStorage) GetBuildingID(ctx context.Context, planetID uuid.UUID, BuildingType consts.BuildingType) (consts.BuildingID, error) {
	const getBuildingIDQuery = `
		SELECT pb.building_id
		FROM session_beta.planet_buildings pb
		JOIN session_beta.s_buildings b ON pb.building_id = b.id
		WHERE pb.planet_id = $1 AND b.building_type = $2;
	`

	var buildingID consts.BuildingID
	err := s.DB.QueryRow(ctx, getBuildingIDQuery, planetID, BuildingType).Scan(&buildingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, models.ErrBuildingNotFound
		}

		return 0, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return buildingID, nil
}
