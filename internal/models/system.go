package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type System struct {
	X consts.PlanetPositionX
	Y consts.PlanetPositionY
}

type SystemPlanets struct {
	System  System
	Planets []PlanetInfo
	NPC     []NPCAttack
}

type PlanetInfo struct {
	ID         uuid.UUID
	Z          consts.PlanetPositionZ
	Type       string
	UserLogin  string
	HasMoon    bool
	Debris     Resources
	AttackedAt time.Time
}

type NPCAttack struct {
	Z          consts.PlanetPositionZ
	AttackedAt time.Time
}
