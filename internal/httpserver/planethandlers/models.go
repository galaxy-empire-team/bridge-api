package planethandlers

import (
	"time"

	"github.com/google/uuid"
)

// Use POST. Change to Get after authorization implementation
type UserIDRequest struct {
	UserID uuid.UUID `json:"userID"`
}

type PlanetResponse struct {
	PlanetID  uuid.UUID               `json:"planetID"`
	X         uint8                   `json:"x"`
	Y         uint8                   `json:"y"`
	Z         uint8                   `json:"z"`
	Resource  PlanetResources         `json:"resources"`
	Buildings map[string]BuildingInfo `json:"buildings"`
	IsCapitol bool                    `json:"isCapitol"`
	HasMoon   bool                    `json:"hasMoon"`
}

type PlanetResources struct {
	Metal   uint64 `json:"metal"`
	Crystal uint64 `json:"crystal"`
	Gas     uint64 `json:"gas"`
}

type ErrorResponse struct {
	Err string `json:"err"`
}

type CreateBuildingRequest struct {
	PlanetID     uuid.UUID `json:"planetID"`
	BuildingType string    `json:"buildingType"`
}

type BuildingInfo struct {
	Level            uint8      `json:"level"`
	MetalPerSecond   uint64     `json:"metalPerSecond"`
	CrystalPerSecond uint64     `json:"crystalPerSecond"`
	GasPerSecond     uint64     `json:"gasPerSecond"`
	Bonuses          *string    `json:"bonuses"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	FinishedAt       *time.Time `json:"finishedAt"`
}
