package systemhandlers

import "github.com/google/uuid"

type SystemPlanetsRequest struct {
	X uint64 `json:"x"`
	Y uint64 `json:"y"`
}

type SystemPlanetsResponse struct {
	X       uint64                `json:"x"`
	Y       uint64                `json:"y"`
	Planets map[uint64]PlanetInfo `json:"planets"`
}

type PlanetInfo struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	UserLogin string    `json:"userLogin"`
	HasMoon   bool      `json:"hasMoon"`
}

type ErrorResponse struct {
	Err string `json:"err"`
}
