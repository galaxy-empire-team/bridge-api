package registry

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type Registry struct {
	buildings     map[consts.BuildingID]BuildingStats
	fleet         map[consts.FleetUnitID]FleetUnitStats
	researches    map[consts.ResearchID]ResearchStats
	missions      map[consts.MissionID]consts.MissionType
	notifications map[consts.NotificationID]consts.NotificationType
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

	researchStats, err := getResearchStats(ctx, connPool)
	if err != nil {
		return nil, fmt.Errorf("getResearchStats(): %w", err)
	}

	missionMapping, err := getMissionMapping(ctx, connPool)
	if err != nil {
		return nil, fmt.Errorf("getMissionMapping(): %w", err)
	}

	notificationMapping, err := getNotificationMapping(ctx, connPool)
	if err != nil {
		return nil, fmt.Errorf("getNotificationMapping(): %w", err)
	}

	return &Registry{
		buildings:     buildingStats,
		fleet:         fleetStats,
		researches:    researchStats,
		missions:      missionMapping,
		notifications: notificationMapping,
	}, nil
}

type buildingBonuses struct {
	FleetBuildSpeed float32 `json:"fleet_build_speed,omitempty"`
	ResearchSpeed   float32 `json:"research_speed,omitempty"`
	BuildSpeed      float32 `json:"building_speed,omitempty"`
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
		FROM session_beta.s_buildings;
	`

	var bonuses buildingBonuses
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
			&bonuses,
			&BuildingStats.UpgradeTimeS,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		BuildingStats.Bonuses = BuildingBonuses{
			FleetBuildSpeed: bonuses.FleetBuildSpeed,
			ResearchSpeed:   bonuses.ResearchSpeed,
			BuildSpeed:      bonuses.BuildSpeed,
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
		FROM session_beta.s_fleet;
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

type researchBonuses struct {
	AvaliableColonizePlanetCount uint8 `json:"availiable_colonize_count"`

	ProductionSpeedImprove   float32 `json:"resource_gain"`
	FleetCostReduce          float32 `json:"fleet_cost_reduce"`
	FleetConstructTimeReduce float32 `json:"fleet_construct_time_reduce"`
	PlanetDefenseImprove     float32 `json:"defense_power"`

	AttackPower          float32 `json:"attack_power"`
	ArmorPower           float32 `json:"armor_strength"`
	AttackOnDefensePower float32 `json:"attack_deffence_power"`
	SpyChanceImprove     float32 `json:"success_spy_chance"`
}

func getResearchStats(ctx context.Context, pool *pgxpool.Pool) (map[consts.ResearchID]ResearchStats, error) {
	const getResearchStatsQuery = `
		SELECT 
			id,
			research_type,
			level,
			metal_cost,
			crystal_cost,
			gas_cost,
			bonuses,
			research_time_s
		FROM session_beta.s_researches;
	`
	var bonuses researchBonuses
	result := make(map[consts.ResearchID]ResearchStats)
	rows, err := pool.Query(ctx, getResearchStatsQuery)
	if err != nil {
		return nil, fmt.Errorf("pool.Query(): %w", err)
	}

	for rows.Next() {
		var researchStats ResearchStats
		err = rows.Scan(
			&researchStats.ID,
			&researchStats.Type,
			&researchStats.Level,
			&researchStats.MetalCost,
			&researchStats.CrystalCost,
			&researchStats.GasCost,
			&bonuses,
			&researchStats.ResearchTimeS,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		researchStats.Bonuses = ResearchBonuses{
			AvaliableColonizePlanetCount: bonuses.AvaliableColonizePlanetCount,
			ProductionSpeedImprove:       bonuses.ProductionSpeedImprove,
			FleetCostReduce:              bonuses.FleetCostReduce,
			FleetConstructTimeReduce:     bonuses.FleetConstructTimeReduce,
			PlanetDefense:                bonuses.PlanetDefenseImprove,
			AttackPower:                  bonuses.AttackPower,
			ArmorPower:                   bonuses.ArmorPower,
			AttackOnDefensePower:         bonuses.AttackOnDefensePower,
			SpyChanceImprove:             bonuses.SpyChanceImprove,
		}

		result[researchStats.ID] = researchStats
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return result, nil
}

func getNotificationMapping(ctx context.Context, pool *pgxpool.Pool) (map[consts.NotificationID]consts.NotificationType, error) {
	const getNotificationMappingQuery = `
		SELECT 
			id,
			notification_type
		FROM session_beta.s_notifications;
	`

	result := make(map[consts.NotificationID]consts.NotificationType)
	rows, err := pool.Query(ctx, getNotificationMappingQuery)
	if err != nil {
		return nil, fmt.Errorf("pool.Query(): %w", err)
	}

	for rows.Next() {
		var id consts.NotificationID
		var notificationType consts.NotificationType
		err = rows.Scan(
			&id,
			&notificationType,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		result[id] = notificationType
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return result, nil
}

func getMissionMapping(ctx context.Context, pool *pgxpool.Pool) (map[consts.MissionID]consts.MissionType, error) {
	const getMissionMappingQuery = `
		SELECT
			id,
			mission_type
		FROM session_beta.s_missions;
	`

	result := make(map[consts.MissionID]consts.MissionType)
	rows, err := pool.Query(ctx, getMissionMappingQuery)
	if err != nil {
		return nil, fmt.Errorf("pool.Query(): %w", err)
	}

	for rows.Next() {
		var id consts.MissionID
		var missionType consts.MissionType
		err = rows.Scan(
			&id,
			&missionType,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		result[id] = missionType
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return result, nil
}
