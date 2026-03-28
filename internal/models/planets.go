package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type Planet struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	Coordinates         Coordinates
	Resources           Resources
	Buildings           []consts.BuildingID
	BuildingsInProgress []BuildingInProgress
	HasMoon             bool
	IsCapitol           bool
	ColonizedAt         time.Time
	UpdatedAt           time.Time
}

type Resources struct {
	Metal     uint64
	Crystal   uint64
	Gas       uint64
	UpdatedAt time.Time
}

type Coordinates struct {
	X consts.PlanetPositionX
	Y consts.PlanetPositionY
	Z consts.PlanetPositionZ
}

type PlanetIDWithCapitol struct {
	PlanetID  uuid.UUID
	IsCapitol bool
}
