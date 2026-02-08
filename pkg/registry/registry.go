package registry

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type Registry struct {
	buildings map[consts.BuildingID]BuildingStats
	fleet     map[consts.FleetUnitID]FleetUnitStats
}

func New(ctx context.Context, connPool *pgxpool.Pool) (*Registry, error) {
	buildingStats, err := getBuildingStats(ctx, connPool)
	if err != nil {
		return nil, fmt.Errorf("getBuildingStats(): %w", err)
	}

	fleetStats, err := getFleetStats(ctx, connPool)
	if err != nil {
		return nil, fmt.Errorf("getFleetStats(): %w", err)
	}

	return &Registry{
		buildings: buildingStats,
		fleet:     fleetStats,
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

func getFleetStats(ctx context.Context, pool *pgxpool.Pool) (map[consts.FleetUnitID]FleetUnitStats, error) {
	const getFleetStatsQuery = `
		SELECT 
			id,
			ship_type,
			attack,
			defense,
			speed,
			cargo_capacity,
			metal_cost,
			crystal_cost,
			gas_cost,
			build_time_s
		FROM session_beta.fleet;
	`

	result := make(map[consts.FleetUnitID]FleetUnitStats)
	rows, err := pool.Query(ctx, getFleetStatsQuery)
	if err != nil {
		return nil, fmt.Errorf("pool.Query(): %w", err)
	}

	for rows.Next() {
		var fleetUnitStats FleetUnitStats
		err = rows.Scan(
			&fleetUnitStats.ID,
			&fleetUnitStats.Type,
			&fleetUnitStats.Attack,
			&fleetUnitStats.Defense,
			&fleetUnitStats.Speed,
			&fleetUnitStats.CargoCapacity,
			&fleetUnitStats.MetalCost,
			&fleetUnitStats.CrystalCost,
			&fleetUnitStats.GasCost,
			&fleetUnitStats.BuildTimeSec,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		result[fleetUnitStats.ID] = fleetUnitStats
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return result, nil
}
