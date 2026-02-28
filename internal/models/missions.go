package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type MissionEvent struct {
	UserID      uuid.UUID
	PlanetFrom  uuid.UUID
	PlanetTo    Coordinates
	Type        consts.MissionID
	Fleet       []PlanetFleetUnitCount
	Cargo       Resources
	IsReturning bool
	StartedAt   time.Time
	FinishedAt  time.Time
}

type UserMission struct {
	Type        string
	PlanetFrom  Coordinates
	PlanetTo    Coordinates
	IsReturning bool
	StartedAt   time.Time
	FinishedAt  time.Time
}
