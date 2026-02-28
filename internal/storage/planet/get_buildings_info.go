package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// GetBuildsInfo retrieves mine infromation from the target planet's buildings
func (s *PlanetStorage) GetBuildingsInfo(ctx context.Context, planetID uuid.UUID, BuildingTypes []consts.BuildingType) (map[consts.BuildingType]models.BuildingInfo, error) {
	const getMineInfoQuery = `
		SELECT 
			b.level,
			b.building_type,
			b.production_s,
			b.bonuses,
			pb.updated_at,
			pb.finished_at
		FROM session_beta.planet_buildings pb
		JOIN session_beta.s_buildings b ON pb.building_id = b.id
		WHERE pb.planet_id = $1 AND b.building_type = ANY($2);
	`

	rows, err := s.DB.Query(ctx, getMineInfoQuery, planetID, BuildingTypes)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}

	var updatedAt time.Time
	var finishedAt *time.Time
	buildingsInfo := make(map[consts.BuildingType]models.BuildingInfo)
	for rows.Next() {
		var buildingInfo models.BuildingInfo
		err = rows.Scan(
			&buildingInfo.Level,
			&buildingInfo.Type,
			&buildingInfo.ProductionS,
			&buildingInfo.Bonuses,
			&updatedAt,
			&finishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		buildingInfo.UpdatedAt = updatedAt.UTC()
		if finishedAt != nil {
			buildingInfo.FinishedAt = finishedAt.UTC()
		}

		buildingsInfo[buildingInfo.Type] = buildingInfo
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return buildingsInfo, nil
}
