package missionhandlers

import (
	"github.com/google/uuid"
)

// Use POST. Change to Get after authorization implementation.
type ColonizeRequest struct {
	UserID     uuid.UUID   `json:"userID"`
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
