package event

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

const (
	baseFleetConstructionSpeedMultiplier = 1.0
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

	speedMultiplier, err := s.calculateSpeedMultiplier(ctx, userID, planet)
	if err != nil {
		return models.FleetUnitConstructionInfo{}, fmt.Errorf("calculateSpeedMultiplier(): %w", err)
	}

	return info, s.txManager.ExecEventTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		fleetConstructEvent, err := s.generateEventForFleetConstruct(ctx, planet, fleet, speedMultiplier, planetRepo)
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

func (s *Service) generateEventForFleetConstruct(ctx context.Context, planetID uuid.UUID, fleet models.FleetUnitCount, speedMultiplier float64, planetRepo TxStorages) (models.FleetConstructEvent, error) {
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
	finishedAt := startedAt.Add(time.Duration(float64(fleetUnitStats.BuildTimeSec*fleet.Count)/speedMultiplier) * time.Second).UTC()
	fleetConstructEvent := models.FleetConstructEvent{
		PlanetID:      planetID,
		FleetID:       fleet.ID,
		Count:         fleet.Count,
		ResourcesCost: resourcesCost,
		StartedAt:     startedAt,
		FinishedAt:    finishedAt,
	}

	return fleetConstructEvent, nil
}

func (s *Service) calculateSpeedMultiplier(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (float64, error) {
	moonInfo, err := s.planetStorage.GetFullMoonInfo(ctx, planetID)
	if err != nil {
		return 0, fmt.Errorf("planetStorage.GetFullMoonInfo(): %w", err)
	}

	speedMultiplier := baseFleetConstructionSpeedMultiplier
	if moonInfo.HasMoon {
		if moonInfo.ActivateUntill.After(time.Now().UTC()) {
			speedMultiplier = baseFleetConstructionSpeedMultiplier + consts.ActiveFleetConstructionMoonSpeedMultiplier
		} else {
			speedMultiplier = baseFleetConstructionSpeedMultiplier + consts.InactiveFleetConstructionMoonSpeedMultiplier
		}
	}

	spaceportStats, err := s.repository.GetBuildingByType(ctx, planetID, consts.BuildingTypeSpaceport)
	if err != nil {
		return 0, fmt.Errorf("repository.GetBuildingByType(): %w", err)
	}

	constructionOptimizationStats, err := s.repository.GetResearchByType(ctx, userID, consts.ResearchTypeConstructionOptimizationTechnology)
	if err != nil {
		return 0, fmt.Errorf("repository.GetResearchByType(): %w", err)
	}

	return float64(spaceportStats.Bonuses.FleetBuildSpeed) * float64(constructionOptimizationStats.Bonuses.FleetConstructTimeReduce) * speedMultiplier, nil
}
