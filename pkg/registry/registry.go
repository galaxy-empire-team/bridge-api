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

	zeroLvlBuildings  map[consts.BuildingType]consts.BuildingID
	zeroLvlResearches map[consts.ResearchType]consts.ResearchID
}

func New(ctx context.Context, connPool *pgxpool.Pool) (*Registry, error) {
	r := &Registry{}

	err := r.fillBuildingStats(ctx, connPool)
	if err != nil {
		return nil, fmt.Errorf("fillBuildingStats(): %w", err)
	}

	err = r.fillFleetStats(ctx, connPool)
	if err != nil {
		return nil, fmt.Errorf("getFleetStats(): %w", err)
	}

	err = r.fillResearchStats(ctx, connPool)
	if err != nil {
		return nil, fmt.Errorf("fillResearchStats(): %w", err)
	}

	err = r.fillMissionMapping(ctx, connPool)
	if err != nil {
		return nil, fmt.Errorf("fillMissionMapping(): %w", err)
	}

	err = r.fillNotificationMapping(ctx, connPool)
	if err != nil {
		return nil, fmt.Errorf("fillNotificationMapping(): %w", err)
	}

	return r, nil
}

type buildingBonuses struct {
	FleetBuildSpeed float32 `json:"fleet_build_speed,omitempty"`
	ResearchSpeed   float32 `json:"research_speed,omitempty"`
	BuildSpeed      float32 `json:"building_speed,omitempty"`
}

func (r *Registry) fillBuildingStats(ctx context.Context, pool *pgxpool.Pool) error {
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
	r.buildings = make(map[consts.BuildingID]BuildingStats)
	r.zeroLvlBuildings = make(map[consts.BuildingType]consts.BuildingID)
	rows, err := pool.Query(ctx, getBuildingStatsQuery)
	if err != nil {
		return fmt.Errorf("pool.Query(): %w", err)
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
			return fmt.Errorf("rows.Scan(): %w", err)
		}

		if BuildingStats.Level == 0 {
			r.zeroLvlBuildings[BuildingStats.Type] = BuildingStats.ID
		}

		BuildingStats.Bonuses = BuildingBonuses{
			FleetBuildSpeed: bonuses.FleetBuildSpeed,
			ResearchSpeed:   bonuses.ResearchSpeed,
			BuildSpeed:      bonuses.BuildSpeed,
		}

		r.buildings[BuildingStats.ID] = BuildingStats
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows.Err(): %w", err)
	}

	return nil
}

func (r *Registry) fillFleetStats(ctx context.Context, pool *pgxpool.Pool) error {
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

	r.fleet = make(map[consts.FleetUnitID]FleetUnitStats)
	rows, err := pool.Query(ctx, getFleetStatsQuery)
	if err != nil {
		return fmt.Errorf("pool.Query(): %w", err)
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
			return fmt.Errorf("rows.Scan(): %w", err)
		}

		r.fleet[fleetUnitStats.ID] = fleetUnitStats
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows.Err(): %w", err)
	}

	return nil
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

func (r *Registry) fillResearchStats(ctx context.Context, pool *pgxpool.Pool) error {
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
	r.researches = make(map[consts.ResearchID]ResearchStats)
	r.zeroLvlResearches = make(map[consts.ResearchType]consts.ResearchID)
	rows, err := pool.Query(ctx, getResearchStatsQuery)
	if err != nil {
		return fmt.Errorf("pool.Query(): %w", err)
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
			return fmt.Errorf("rows.Scan(): %w", err)
		}

		if researchStats.Level == 0 {
			r.zeroLvlResearches[researchStats.Type] = researchStats.ID
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

		r.researches[researchStats.ID] = researchStats
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows.Err(): %w", err)
	}

	return nil
}

func (r *Registry) fillNotificationMapping(ctx context.Context, pool *pgxpool.Pool) error {
	const getNotificationMappingQuery = `
		SELECT 
			id,
			notification_type
		FROM session_beta.s_notifications;
	`

	r.notifications = make(map[consts.NotificationID]consts.NotificationType)
	rows, err := pool.Query(ctx, getNotificationMappingQuery)
	if err != nil {
		return fmt.Errorf("pool.Query(): %w", err)
	}

	for rows.Next() {
		var id consts.NotificationID
		var notificationType consts.NotificationType
		err = rows.Scan(
			&id,
			&notificationType,
		)
		if err != nil {
			return fmt.Errorf("rows.Scan(): %w", err)
		}

		r.notifications[id] = notificationType
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows.Err(): %w", err)
	}

	return nil
}

func (r *Registry) fillMissionMapping(ctx context.Context, pool *pgxpool.Pool) error {
	const getMissionMappingQuery = `
		SELECT
			id,
			mission_type
		FROM session_beta.s_missions;
	`

	r.missions = make(map[consts.MissionID]consts.MissionType)
	rows, err := pool.Query(ctx, getMissionMappingQuery)
	if err != nil {
		return fmt.Errorf("pool.Query(): %w", err)
	}

	for rows.Next() {
		var id consts.MissionID
		var missionType consts.MissionType
		err = rows.Scan(
			&id,
			&missionType,
		)
		if err != nil {
			return fmt.Errorf("rows.Scan(): %w", err)
		}

		r.missions[id] = missionType
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows.Err(): %w", err)
	}

	return nil
}
