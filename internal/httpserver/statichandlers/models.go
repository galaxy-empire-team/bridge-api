package statichandlers

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type ErrorResponse struct {
	Err string `json:"err"`
}

type BuildingStats struct {
	ID           consts.BuildingID    `json:"id"`
	Type         consts.BuildingType  `json:"type"`
	Level        consts.BuildingLevel `json:"level"`
	MetalCost    uint64               `json:"metalCost"`
	CrystalCost  uint64               `json:"crystalCost"`
	GasCost      uint64               `json:"gasCost"`
	ProductionS  uint64               `json:"productionS"`
	Bonuses      BuildingBonuses      `json:"bonuses"`
	UpgradeTimeS uint64               `json:"upgradeTimeS"`
}

type BuildingBonuses struct {
	FleetBuildSpeed float32 `json:"fleet_build_speed,omitempty"`
	ResearchSpeed   float32 `json:"research_speed,omitempty"`
	BuildSpeed      float32 `json:"building_speed,omitempty"`
}
