package event

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

	if err := s.repository.CheckPlanetOwner(ctx, userID, planet); err != nil {
		return models.FleetUnitConstructionInfo{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	fleetConstructionInProgress, err := s.eventStorage.CheckFleetConstruction(ctx, planet)
	if err != nil {
		return models.FleetUnitConstructionInfo{}, fmt.Errorf("planetStorage.CheckFleetConstruction(): %w", err)
	}
	if fleetConstructionInProgress {
		return models.FleetUnitConstructionInfo{}, models.ErrFleetConstructionInProgress
	}

	err = s.repository.RecalcResources(ctx, userID, planet)
	if err != nil {
		return models.FleetUnitConstructionInfo{}, fmt.Errorf("recalcResources(): %w", err)
	}

	info := models.FleetUnitConstructionInfo{
		FleetID: fleet.ID,
		Count:   fleet.Count,
	}

	return info, s.txManager.ExecEventTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
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

	resourcesCost := models.Resources{
		Metal:   fleetUnitStats.MetalCost * fleet.Count,
		Crystal: fleetUnitStats.CrystalCost * fleet.Count,
		Gas:     fleetUnitStats.GasCost * fleet.Count,
	}

	if resources.Metal < resourcesCost.Metal ||
		resources.Crystal < resourcesCost.Crystal ||
		resources.Gas < resourcesCost.Gas {
		return models.FleetConstructEvent{}, models.ErrNotEnoughResources
	}

	leftResources := models.Resources{
		Metal:     resources.Metal - resourcesCost.Metal,
		Crystal:   resources.Crystal - resourcesCost.Crystal,
		Gas:       resources.Gas - resourcesCost.Gas,
		UpdatedAt: resources.UpdatedAt,
	}

	err = planetRepo.SetResources(ctx, planetID, leftResources)
	if err != nil {
		return models.FleetConstructEvent{}, fmt.Errorf("planetRepo.SetResources(): %w", err)
	}

	startedAt := time.Now().UTC()
	fleetConstructEvent := models.FleetConstructEvent{
		PlanetID:      planetID,
		FleetID:       fleet.ID,
		Count:         fleet.Count,
		ResourcesCost: resourcesCost,
		StartedAt:     startedAt,
		FinishedAt:    startedAt.Add(time.Duration(fleetUnitStats.BuildTimeSec*fleet.Count) * time.Second).UTC(),
	}

	return fleetConstructEvent, nil
}
