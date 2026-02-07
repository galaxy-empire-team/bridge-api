package registry

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type Registry struct {
	buildings map[consts.BuildingID]BuildingStats
}

func New(ctx context.Context, connPool *pgxpool.Pool) (*Registry, error) {
	BuildingStats, err := getBuildingStats(ctx, connPool)
	if err != nil {
		return nil, fmt.Errorf("getBuildingStats(): %w", err)
	}

	return &Registry{
		buildings: BuildingStats,
	}, nil
}

func getBuildingStats(ctx context.Context, pool *pgxpool.Pool) (map[consts.BuildingID]BuildingStats, error) {
	const getBuildingStatsQuery = `
		SELECT 
			id,
			building_type, 
			level, 
			metal_cost,
			crystal_cost,
			gas_cost,
			production_s,
			bonuses,
			upgrade_time_s
		FROM session_beta.buildings;
	`

	result := make(map[consts.BuildingID]BuildingStats)
	rows, err := pool.Query(ctx, getBuildingStatsQuery)
	if err != nil {
		return nil, fmt.Errorf("pool.Query(): %w", err)
	}

	for rows.Next() {
		var BuildingStats BuildingStats
		err = rows.Scan(
			&BuildingStats.ID,
			&BuildingStats.Type,
			&BuildingStats.Level,
			&BuildingStats.MetalCost,
			&BuildingStats.CrystalCost,
			&BuildingStats.GasCost,
			&BuildingStats.ProductionS,
			&BuildingStats.Bonuses,
			&BuildingStats.UpgradeTimeS,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		result[BuildingStats.ID] = BuildingStats
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return result, nil
}
