package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type MissionEvent struct {
	UserID     uuid.UUID
	PlanetFrom uuid.UUID
	PlanetTo   Coordinates
	Type       consts.MissionType
	StartedAt  time.Time
	FinishedAt time.Time
}
