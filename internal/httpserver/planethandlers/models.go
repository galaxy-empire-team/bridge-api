package planethandlers


import (	
	"github.com/google/uuid"
)
// Use POST. Change to Get after authorization implementation
type UserRequest struct {
	UserID uuid.UUID `json:"userID"`
}

type PlanetResponse struct {
	PlanetID uuid.UUID `json:"planetID"`
	X 	  uint8    `json:"x"`
	Y 	  uint8    `json:"y"`
	Z 	  uint8    `json:"z"`
	Resource PlanetResources `json:"resources"`
	IsCapitol bool     `json:"isCapitol"`
	HasMoon  bool     `json:"hasMoon"`
}

type PlanetResources struct {
	Metal   uint64 `json:"metal"`
	Crystal uint64 `json:"crystal"`
	Gas     uint64 `json:"gas"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
