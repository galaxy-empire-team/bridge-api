package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) StartFleetConstruction(ctx context.Context, userID uuid.UUID, planet uuid.UUID, fleet models.FleetUnitCount) (models.FleetUnitConstructionInfo, error) {
	if fleet.ID == 0 || fleet.Count == 0 {
		return models.FleetUnitConstructionInfo{}, models.ErrInvalidFleetConstructionRequest
	}

	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planet)
	if err != nil {
		return models.FleetUnitConstructionInfo{}, fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.FleetUnitConstructionInfo{}, models.ErrPlanetDoesNotBelongToUser
	}

	fleetConstructionInProgress, err := s.planetStorage.CheckFleetConstruction(ctx, planet)
	if err != nil {
		return models.FleetUnitConstructionInfo{}, fmt.Errorf("planetStorage.CheckFleetConstruction(): %w", err)
	}
	if fleetConstructionInProgress {
		return models.FleetUnitConstructionInfo{}, models.ErrFleetConstructionInProgress
	}

	err = s.resourceRepository.RecalcResources(ctx, userID, planet)
	if err != nil {
		return models.FleetUnitConstructionInfo{}, fmt.Errorf("recalcResources(): %w", err)
	}

	info := models.FleetUnitConstructionInfo{
		FleetID: fleet.ID,
		Count:   fleet.Count,
	}

	return info, s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		fleetConstructEvent, err := s.generateEventForFleetConstruct(ctx, planet, fleet, planetRepo)
		if err != nil {
			return fmt.Errorf("generateEventForFleetConstruct(): %w", err)
		}

		err = planetRepo.CreateFleetConstructEvent(ctx, fleetConstructEvent)
		if err != nil {
			return fmt.Errorf("planetStorage.CreateFleetConstructEvent(): %w", err)
		}

		info.StartedAt = fleetConstructEvent.StartedAt
		info.FinishedAt = fleetConstructEvent.FinishedAt

		return nil
	})
}

func (s *Service) generateEventForFleetConstruct(ctx context.Context, planetID uuid.UUID, fleet models.FleetUnitCount, planetRepo TxStorages) (models.FleetConstructEvent, error) {
	fleetUnitStats, err := s.registry.GetFleetUnitStatsByID(fleet.ID)
	if err != nil {
		return models.FleetConstructEvent{}, fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
	}

	// Calculate resources
	resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return models.FleetConstructEvent{}, fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
	}

	if resources.Metal < fleetUnitStats.MetalCost*fleet.Count ||
		resources.Crystal < fleetUnitStats.CrystalCost*fleet.Count ||
		resources.Gas < fleetUnitStats.GasCost*fleet.Count {
		return models.FleetConstructEvent{}, models.ErrNotEnoughResources
	}

	leftResources := models.Resources{
		Metal:     resources.Metal - (fleetUnitStats.MetalCost * fleet.Count),
		Crystal:   resources.Crystal - (fleetUnitStats.CrystalCost * fleet.Count),
		Gas:       resources.Gas - (fleetUnitStats.GasCost * fleet.Count),
		UpdatedAt: resources.UpdatedAt,
	}

	err = planetRepo.SetResources(ctx, planetID, leftResources)
	if err != nil {
		return models.FleetConstructEvent{}, fmt.Errorf("planetRepo.SetResources(): %w", err)
	}

	startedAt := time.Now().UTC()
	fleetConstructEvent := models.FleetConstructEvent{
		PlanetID:   planetID,
		FleetID:    fleet.ID,
		Count:      fleet.Count,
		StartedAt:  startedAt,
		FinishedAt: startedAt.Add(time.Duration(fleetUnitStats.BuildTimeSec*fleet.Count) * time.Second).UTC(),
	}

	return fleetConstructEvent, nil
}
