package planet

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) StartBuildingUpgrade(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID) (models.FinishTime, error) {
	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetID)
	if err != nil {
		return models.FinishTime{}, fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.FinishTime{}, models.ErrPlanetDoesNotBelongToUser
	}

	currentBuildsCount, err := s.planetStorage.GetBuildsInProgressCount(ctx, planetID)
	if err != nil {
		return models.FinishTime{}, fmt.Errorf("planetStorage.GetBuildsInProgressCount(): %w", err)
	}

	if currentBuildsCount >= consts.MaxBuildingsInProgress {
		return models.FinishTime{}, models.ErrTooManyBuildingsInProgress
	}

	planetBuildingIDs, err := s.planetStorage.GetAllPlanetBuildings(ctx, planetID)
	if err != nil {
		return models.FinishTime{}, fmt.Errorf("planetStorage.GetAllPlanetBuildings(): %w", err)
	}

	if !slices.Contains(planetBuildingIDs, buildingID) {
		buildStats, err := s.registry.GetBuildingStatsByID(buildingID)
		if err != nil {
			return models.FinishTime{}, fmt.Errorf("registry.GetBuildingStatsByID(): %w", err)
		}

		if buildStats.Level != consts.ZeroBuildingLevel {
			return models.FinishTime{}, models.ErrBuildingNotFound
		}
	}

	err = s.recalcResources(ctx, userID, planetID)
	if err != nil {
		return models.FinishTime{}, fmt.Errorf("recalcResources(): %w", err)
	}

	FinishTime := models.FinishTime{}

	return FinishTime, s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		buildEvent, err := s.generateEventForExistingBuilding(ctx, planetID, buildingID, planetRepo)
		if err != nil {
			return fmt.Errorf("generateEventForExistingBuilding(): %w", err)
		}

		err = planetRepo.CreateBuildingEvent(ctx, buildEvent)
		if err != nil {
			return fmt.Errorf("planetStorage.CreateBuildingEvent(): %w", err)
		}

		FinishTime.StartedAt = buildEvent.StartedAt
		FinishTime.FinishedAt = buildEvent.FinishedAt

		return nil
	})
}

func (s *Service) generateEventForExistingBuilding(ctx context.Context, planetID uuid.UUID, buildingID consts.BuildingID, planetRepo TxStorages) (models.BuildEvent, error) {
	nextLvlBuildingID, err := s.registry.GetBuildingNextLvlID(buildingID)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("registry.GetBuildingNextLvlID(): %w", err)
	}

	nextLvlStats, err := s.registry.GetBuildingStatsByID(nextLvlBuildingID)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("registry.GetBuildingStatsByID(): %w", err)
	}

	resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
	}

	if resources.Metal < nextLvlStats.MetalCost ||
		resources.Crystal < nextLvlStats.CrystalCost ||
		resources.Gas < nextLvlStats.GasCost {
		return models.BuildEvent{}, models.ErrNotEnoughResources
	}

	leftResources := models.Resources{
		Metal:     resources.Metal - nextLvlStats.MetalCost,
		Crystal:   resources.Crystal - nextLvlStats.CrystalCost,
		Gas:       resources.Gas - nextLvlStats.GasCost,
		UpdatedAt: resources.UpdatedAt,
	}

	err = planetRepo.SetResources(ctx, planetID, leftResources)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("planetRepo.SetResources(): %w", err)
	}

	startedAt := time.Now().UTC()
	buildEvent := models.BuildEvent{
		PlanetID:   planetID,
		BuildingID: buildingID,
		StartedAt:  startedAt,
		FinishedAt: startedAt.Add(time.Duration(nextLvlStats.UpgradeTimeS) * time.Second).UTC(),
	}

	return buildEvent, nil
}
