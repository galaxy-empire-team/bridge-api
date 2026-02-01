package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type BuildEvent struct {
	PlanetID     uuid.UUID
	BuildingType consts.BuildingType
	StartedAt    time.Time
	FinishedAt   time.Time
}

type BuildingStats struct {
	Level        uint8
	Type         consts.BuildingType
	MetalCost    uint64
	CrystalCost  uint64
	GasCost      uint64
	ProductionS  uint64
	Bonuses      *string
	UpgradeTimeS uint16
}

type BuildingInfo struct {
	Level       uint8
	Type        consts.BuildingType
	ProductionS uint64
	Bonuses     *string
	UpdatedAt   time.Time
	FinishedAt  time.Time
}
