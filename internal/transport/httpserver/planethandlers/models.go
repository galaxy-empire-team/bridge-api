package planethandlers

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type ErrorResponse struct {
	Err string `json:"err"`
}

type PlanetIDRequest struct {
	PlanetID uuid.UUID `json:"planetID"`
}

type CapitolIDResponse struct {
	CapitolPlanetID uuid.UUID `json:"capitolID"`
}

type StartBuildingUpgradeRequest struct {
	PlanetID   uuid.UUID         `json:"planetID"`
	BuildingID consts.BuildingID `json:"buildingID"`
}

type CancelBuildingUpgradeRequest struct {
	PlanetID   uuid.UUID         `json:"planetID"`
	BuildingID consts.BuildingID `json:"buildingID"`
}

type StartResearchRequest struct {
	PlanetID   uuid.UUID         `json:"planetID"`
	ResearchID consts.ResearchID `json:"researchID"`
}

type CancelResearchRequest struct {
	PlanetID   uuid.UUID         `json:"planetID"`
	ResearchID consts.ResearchID `json:"researchID"`
}

type StartFleetConstructionRequest struct {
	PlanetID uuid.UUID          `json:"planetID"`
	FleetID  consts.FleetUnitID `json:"fleetID"`
	Count    uint64             `json:"count"`
}

type CancelFleetConstructionRequest struct {
	PlanetID uuid.UUID `json:"planetID"`
}

type PlanetResponse struct {
	PlanetID    uuid.UUID       `json:"planetID"`
	X           uint8           `json:"x"`
	Y           uint16          `json:"y"`
	Z           uint8           `json:"z"`
	Resources   PlanetResources `json:"resources"`
	IsCapitol   bool            `json:"isCapitol"`
	HasMoon     bool            `json:"hasMoon"`
	ColonizedAt time.Time       `json:"colonizedAt"`
}

type BuildingInProgress struct {
	BuildingID consts.BuildingID `json:"id"`
	StartedAt  time.Time         `json:"startedAt"`
	FinishedAt time.Time         `json:"finishedAt"`
}

type PlanetBuildingsResponse struct {
	Buildings          []consts.BuildingID  `json:"buildings"`
	BuildingInProgress []BuildingInProgress `json:"progress,omitempty"`
}

type PlanetResources struct {
	Metal     uint64    `json:"metal"`
	Crystal   uint64    `json:"crystal"`
	Gas       uint64    `json:"gas"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserPlanetsResponse struct {
	Planets []ShortPlanet `json:"planets"`
}

type ShortPlanet struct {
	PlanetID    uuid.UUID `json:"planetID"`
	X           uint8     `json:"x"`
	Y           uint16    `json:"y"`
	Z           uint8     `json:"z"`
	IsCapitol   bool      `json:"isCapitol"`
	HasMoon     bool      `json:"hasMoon"`
	ColonizedAt time.Time `json:"colonizedAt"`
}

type FleetPlanetsResponse struct {
	Fleet             []FleetUnitCount          `json:"fleet"`
	FleetConstruction FleetConstructionResponse `json:"construction,omitzero"`
}

type FleetUnitCount struct {
	ID    uint8  `json:"id"`
	Count uint64 `json:"count"`
}

type FinishTimeResponse struct {
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`
}

type FleetConstructionResponse struct {
	FleetID    consts.FleetUnitID `json:"id"`
	Count      uint64             `json:"count"`
	StartedAt  time.Time          `json:"startedAt"`
	FinishedAt time.Time          `json:"finishedAt"`
}

type ResearchesResponse struct {
	Researches       []consts.ResearchID        `json:"researches"`
	ResearchProgress []ResearchProgressResponse `json:"progress"`
}

type ResearchProgressResponse struct {
	ResearchID consts.ResearchID `json:"id"`
	StartedAt  time.Time         `json:"startedAt"`
	FinishedAt time.Time         `json:"finishedAt"`
}

type UserResourcesResponse struct {
	UserResources UserResources `json:"resources"`
	Boosts        []UserBoost   `json:"boosts"`
}

type UserResources struct {
	Matter uint64 `json:"matter"`
	Doreye uint64 `json:"doreye"`
}

type UserBoost struct {
	ID    consts.BoostID `json:"id"`
	Count uint64         `json:"count"`
}
