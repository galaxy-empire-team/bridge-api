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
