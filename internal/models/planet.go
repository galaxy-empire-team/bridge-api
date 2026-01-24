package models

import (
	"time"

	"github.com/google/uuid"
)

type Planet struct {
	ID          uuid.UUID
	Location    Location
	Resources   Resources
	Buildings   map[BuildingType]BuildingInfo
	HasMoon     bool
	IsCapitol   bool
	ColonizedAt time.Time
}

type Resources struct {
	Metal     uint64
	Crystal   uint64
	Gas       uint64
	UpdatedAt time.Time
}

type Location struct {
	X uint8
	Y uint8
	Z uint8
}

type PlanetIDWithCapitol struct {
	PlanetID  uuid.UUID
	IsCapitol bool
}
