package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type BuildEvent struct {
	PlanetID   uuid.UUID
	BuildingID consts.BuildingID
	StartedAt  time.Time
	FinishedAt time.Time
}

type BuildingStats struct {
	Level        consts.BuildingLevel
	Type         consts.BuildingType
	MetalCost    uint64
	CrystalCost  uint64
	GasCost      uint64
	ProductionS  uint64
	Bonuses      map[string]float64
	UpgradeTimeS uint64
}

type BuildingInfo struct {
	Level       consts.BuildingLevel
	Type        consts.BuildingType
	ProductionS uint64
	Bonuses     map[string]float64
	UpdatedAt   time.Time
	FinishedAt  time.Time
}
