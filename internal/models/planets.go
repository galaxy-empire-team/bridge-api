package models

import (
	"time"

	"github.com/google/uuid"
)

type Planet struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Coordinates Coordinates
	Resources   Resources
	HasMoon     bool
	IsCapitol   bool
	ColonizedAt time.Time
	UpdatedAt   time.Time
}

type PlanetIDCapitol struct {
	PlanetID  uuid.UUID
	IsCapitol bool
}
