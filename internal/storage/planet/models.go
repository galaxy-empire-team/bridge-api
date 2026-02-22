package planet

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type planetToColonize struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Coordinates coordinates
	HasMoon     bool
	IsCapitol   bool
}

type coordinates struct {
	X consts.PlanetPositionX
	Y consts.PlanetPositionY
	Z consts.PlanetPositionZ
}

type finishedBuilding struct {
	Type       string
	Level      uint8
	UpdatedAt  time.Time
	FinishedAt time.Time
}
