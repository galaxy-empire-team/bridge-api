package systemhandlers

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type SystemPlanetsRequest struct {
	X consts.PlanetPositionX `json:"x"`
	Y consts.PlanetPositionY `json:"y"`
}

type SystemPlanetsResponse struct {
	X       consts.PlanetPositionX                `json:"x"`
	Y       consts.PlanetPositionY                `json:"y"`
	Planets map[consts.PlanetPositionZ]PlanetInfo `json:"planets"`
	NPC     []NPC                                 `json:"npc"`
}

type PlanetInfo struct {
	ID         uuid.UUID `json:"id"`
	Type       string    `json:"type"`
	UserLogin  string    `json:"userLogin"`
	HasMoon    bool      `json:"hasMoon"`
	Debris     Debris    `json:"debris"`
	AttackedAt time.Time `json:"attackedAt,omitzero"`
}

type NPC struct {
	Z          consts.PlanetPositionZ `json:"z"`
	AttackedAt time.Time              `json:"attackedAt"`
}

type ErrorResponse struct {
	Err string `json:"err"`
}

type Debris struct {
	Metal   uint64 `json:"metal"`
	Crystal uint64 `json:"crystal"`
}
