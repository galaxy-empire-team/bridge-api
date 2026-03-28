package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type FleetConstructEvent struct {
	PlanetID   uuid.UUID
	FleetID    consts.FleetUnitID
	Count      uint64
	StartedAt  time.Time
	FinishedAt time.Time
}

type FleetUnitCount struct {
	ID    consts.FleetUnitID
	Count uint64
}
