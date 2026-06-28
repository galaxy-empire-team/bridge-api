package eventhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type EventService interface {
	StartBuildingUpgrade(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID) (models.FinishTime, error)
	StartResearch(ctx context.Context, userID uuid.UUID, currentPlanet uuid.UUID, currentResearchID consts.ResearchID) (models.ResearchProgressInfo, error)
	StartFleetConstruction(ctx context.Context, userID uuid.UUID, currentPlanet uuid.UUID, fleet models.FleetUnitCount) (models.FleetUnitConstructionInfo, error)
	CancelBuildingUpgrade(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID) error
	CancelResearch(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, researchID consts.ResearchID) error
	CancelFleetConstruction(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error
	BoostResearch(ctx context.Context, userID uuid.UUID, researchID consts.ResearchID, boost models.UserBoost) (models.EventFinishTime, error)
	BoostBuildingUpgrade(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID, boost models.UserBoost) (models.EventFinishTime, error)
	BoostFleetConstruction(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, boost models.UserBoost) (models.EventFinishTime, error)
}
