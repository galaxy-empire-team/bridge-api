package registry

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type BuildingStats struct {
	ID           consts.BuildingID
	Type         consts.BuildingType
	Level        consts.BuildingLevel
	MetalCost    uint64
	CrystalCost  uint64
	GasCost      uint64
	ProductionS  uint64
	Bonuses      BuildingBonuses
	UpgradeTimeS uint64
}

type BuildingBonuses struct {
	FleetBuildSpeed float32
	ResearchSpeed   float32
	BuildSpeed      float32
}

type FleetUnitStats struct {
	ID            consts.FleetUnitID
	Type          consts.FleetUnitType
	Attack        uint64
	Defense       uint64
	Speed         uint64
	MetalCost     uint64
	CrystalCost   uint64
	GasCost       uint64
	CargoCapacity uint64
	BuildTimeSec  uint64
}

type ResearchStats struct {
	ID            consts.ResearchID
	Type          consts.ResearchType
	Level         consts.ResearchLevel
	MetalCost     uint64
	CrystalCost   uint64
	GasCost       uint64
	Bonuses       ResearchBonuses
	ResearchTimeS uint64
}

type ResearchBonuses struct {
	AvaliableColonizePlanetCount uint8
	ProductionSpeedMuliplier     float32
	FleetCostReduce              float32
	FleetConstructTimeReduce     float32

	AttackPower         float32
	ArmorPower          float32
	LootingNPCMuliplier float32
	SpyChanceMuliplier  float32
}

type Resources struct {
	Metal   uint64
	Crystal uint64
	Gas     uint64
}

type FleetUnitCount struct {
	ID    consts.FleetUnitID
	Count uint64
}

type NPCStats struct {
	Tier       uint8
	Name       string
	PositionZ  consts.PlanetPositionZ
	Researches []consts.ResearchID
	Resources  Resources
	Fleet      []FleetUnitCount
	LootFleet  []FleetUnitCount
}

type BoostStats struct {
	ID         consts.BoostID
	Tier       consts.BoostTier
	BoostTimeS uint64
}
