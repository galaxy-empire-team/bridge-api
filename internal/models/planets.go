package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type Planet struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Coordinates Coordinates
	Resources   Resources
	Buildings   map[consts.BuildingType]BuildingInfo
	HasMoon     bool
	IsCapitol   bool
	ColonizedAt time.Time
	UpdatedAt   time.Time
}

type Resources struct {
	Metal     uint64
	Crystal   uint64
	Gas       uint64
	UpdatedAt time.Time
}

type Coordinates struct {
	X uint8
	Y uint8
	Z uint8
}

type PlanetIDWithCapitol struct {
	PlanetID  uuid.UUID
	IsCapitol bool
}
