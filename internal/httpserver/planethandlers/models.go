package planethandlers

import (
	"time"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Err string `json:"err"`
}

type PlanetIDRequest struct {
	PlanetID uuid.UUID `json:"planetID"`
}

type UpgradeBuildingRequest struct {
	PlanetID   uuid.UUID         `json:"planetID"`
	BuildingID consts.BuildingID `json:"buildingID"`
}

type PlanetResponse struct {
	PlanetID            uuid.UUID            `json:"planetID"`
	X                   uint8                `json:"x"`
	Y                   uint16               `json:"y"`
	Z                   uint8                `json:"z"`
	Resources           PlanetResources      `json:"resources"`
	BuildingIDs         []uint16             `json:"buildings"`
	BuildingsInProgress []BuildingInProgress `json:"buildingsInProgress,omitempty"`
	IsCapitol           bool                 `json:"isCapitol"`
	HasMoon             bool                 `json:"hasMoon"`
}

type BuildingInProgress struct {
	BuildingID consts.BuildingID `json:"buildingID"`
	StartedAt  time.Time         `json:"startedAt"`
	FinishedAt time.Time         `json:"finishedAt"`
}

type PlanetResources struct {
	Metal   uint64 `json:"metal"`
	Crystal uint64 `json:"crystal"`
	Gas     uint64 `json:"gas"`
}

type UserPlanetsResponse struct {
	Planets []ShortPlanet `json:"planets"`
}

type ShortPlanet struct {
	PlanetID    uuid.UUID       `json:"planetID"`
	X           uint8           `json:"x"`
	Y           uint16          `json:"y"`
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
