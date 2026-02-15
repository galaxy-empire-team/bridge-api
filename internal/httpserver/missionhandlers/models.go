package missionhandlers

import (
	"time"

	"github.com/google/uuid"
)

type ColonizeRequest struct {
	PlanetFrom uuid.UUID   `json:"planetFrom"`
	PlanetTo   Coordinates `json:"planetTo"`
}

type Coordinates struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
	Z uint8 `json:"z"`
}

type ErrorResponse struct {
	Err string `json:"err"`
}

type AttackRequest struct {
	PlanetFrom     uuid.UUID        `json:"planetFrom"`
	PlanetTo       Coordinates      `json:"planetTo"`
	FleetUnitCount []FleetUnitCount `json:"fleet"`
}

type FleetUnitCount struct {
	ID    uint8  `json:"id"`
	Count uint64 `json:"count"`
}

type UserMissionsResponse struct {
	Missions []Mission `json:"missions"`
}

type Mission struct {
	Type        string      `json:"type"`
	PlanetFrom  Coordinates `json:"planetFrom"`
	PlanetTo    Coordinates `json:"planetTo"`
	IsReturning bool        `json:"isReturning"`
	StartedAt   time.Time   `json:"startedAt"`
	FinishedAt  time.Time   `json:"finishedAt"`
}
