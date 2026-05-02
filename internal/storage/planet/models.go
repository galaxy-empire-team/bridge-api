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
	Resources   resources
	HasMoon     bool
	IsCapitol   bool
}

type coordinates struct {
	X consts.PlanetPositionX
	Y consts.PlanetPositionY
	Z consts.PlanetPositionZ
}

type resources struct {
	Metal   uint64
	Crystal uint64
	Gas     uint64
}

type finishedBuilding struct {
	Type       string
	Level      uint8
	UpdatedAt  time.Time
	FinishedAt time.Time
}
