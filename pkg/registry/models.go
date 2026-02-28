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
	Bonuses      map[string]float64
	UpgradeTimeS uint64
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

	ProductionSpeedImprove   float32
	FleetCostReduce          float32
	FleetConstructTimeReduce float32
	PlanetDefense            float32

	AttackPower          float32
	ArmorPower           float32
	AttackOnDefensePower float32
	SpyChanceImprove     float32
}
