package eventhandlers

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type ErrorResponse struct {
	Err string `json:"err"`
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

type BoostResearchRequest struct {
	ResearchID consts.ResearchID `json:"researchID"`
	Boost      UserBoost         `json:"boost"`
}

type BoostBuildingUpgradeRequest struct {
	PlanetID   uuid.UUID         `json:"planetID"`
	BuildingID consts.BuildingID `json:"buildingID"`
	Boost      UserBoost         `json:"boost"`
}

type BoostFleetConstructionRequest struct {
	PlanetID uuid.UUID `json:"planetID"`
	Boost    UserBoost `json:"boost"`
}

type StartFleetConstructionRequest struct {
	PlanetID uuid.UUID          `json:"planetID"`
	FleetID  consts.FleetUnitID `json:"fleetID"`
	Count    uint64             `json:"count"`
}

type CancelFleetConstructionRequest struct {
	PlanetID uuid.UUID `json:"planetID"`
}

type FinishTimeResponse struct {
	StartedAt  time.Time `json:"startedAt,omitzero"`
	FinishedAt time.Time `json:"finishedAt"`
}

type FleetConstructionResponse struct {
	FleetID    consts.FleetUnitID `json:"id"`
	Count      uint64             `json:"count"`
	StartedAt  time.Time          `json:"startedAt"`
	FinishedAt time.Time          `json:"finishedAt"`
}

type ResearchProgressResponse struct {
	ResearchID consts.ResearchID `json:"id"`
	StartedAt  time.Time         `json:"startedAt"`
	FinishedAt time.Time         `json:"finishedAt"`
}

type UserBoost struct {
	ID    consts.BoostID `json:"id"`
	Count uint64         `json:"count"`
}
