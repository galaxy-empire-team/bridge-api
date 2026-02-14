package planethandlers

import (
	"time"

	"github.com/google/uuid"
)

// Use POST. Change to Get after authorization implementation.
type UserIDRequest struct {
	UserID uuid.UUID `json:"userID"`
}

type PlanetIDRequest struct {
	PlanetID uuid.UUID `json:"planetID"`
}

type UpgradeBuildingRequest struct {
	PlanetID     uuid.UUID `json:"planetID"`
	BuildingType string    `json:"buildingType"`
}

type GetBuildStatsRequest struct {
	BuildingType string `json:"buildingType"`
	Level        uint8  `json:"level"`
}

type GetPlanetResponse struct {
	PlanetID  uuid.UUID               `json:"planetID"`
	X         uint8                   `json:"x"`
	Y         uint8                   `json:"y"`
	Z         uint8                   `json:"z"`
	Resources PlanetResources         `json:"resources"`
	Buildings map[string]BuildingInfo `json:"buildings"`
	IsCapitol bool                    `json:"isCapitol"`
	HasMoon   bool                    `json:"hasMoon"`
}

type GetBuildStatsResponse struct {
	Type                 string             `json:"type"`
	Level                uint8              `json:"level"`
	MetalCost            uint64             `json:"metalCost"`
	CrystalCost          uint64             `json:"crystalCost"`
	GasCost              uint64             `json:"gasCost"`
	ProductionS          uint64             `json:"productionPerSecond"`
	Bonuses              map[string]float64 `json:"bonuses,omitempty"`
	UpgradeTimeInSeconds uint64             `json:"upgradeTimeInSeconds"`
}

type ErrorResponse struct {
	Err string `json:"err"`
}

type BuildingInfo struct {
	Level       uint8              `json:"level"`
	ProductionS uint64             `json:"productionPerSecond"`
	Bonuses     map[string]float64 `json:"bonuses,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty"`
	FinishedAt  time.Time          `json:"finishedAt,omitempty"`
}

type PlanetResources struct {
	Metal   uint64 `json:"metal"`
	Crystal uint64 `json:"crystal"`
	Gas     uint64 `json:"gas"`
}

type UserPlanetsResponse struct {
	Planets []GetShortPlanet `json:"planets"`
}

type GetShortPlanet struct {
	PlanetID    uuid.UUID       `json:"planetID"`
	X           uint8           `json:"x"`
	Y           uint8           `json:"y"`
	Z           uint8           `json:"z"`
	Resources   PlanetResources `json:"resources"`
	IsCapitol   bool            `json:"isCapitol"`
	HasMoon     bool            `json:"hasMoon"`
	ColonizedAt time.Time       `json:"colonizedAt"`
}

type FleetPlanetsResponse struct {
	Fleet []FleetUnitCount `json:"fleet"`
}

type FleetUnitCount struct {
	ID    uint8  `json:"id"`
	Count uint64 `json:"count"`
}
