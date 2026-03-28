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

func (s *Service) StartResearch(ctx context.Context, userID uuid.UUID, planet uuid.UUID, currentResearchID consts.ResearchID) error {
	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planet)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.ErrPlanetDoesNotBelongToUser
	}

	researchInProgress, err := s.researchStorage.CheckResearchInProgress(ctx, userID)
	if err != nil {
		return fmt.Errorf("researchStorage.CheckResearchInProgress(): %w", err)
	}
	if researchInProgress {
		return models.ErrResearchInProgress
	}

	userResearches, err := s.researchStorage.GetUserResearches(ctx, userID)
	if err != nil {
		return fmt.Errorf("researchStorage.GetUserResearches(): %w", err)
	}
	if !slices.Contains(userResearches, currentResearchID) {
		return models.ErrUserHasNotResearch
	}

	err = s.recalcResources(ctx, userID, planet)
	if err != nil {
		return fmt.Errorf("recalcResources(): %w", err)
	}

	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		researchEvent, err := s.generateEventForResearch(ctx, userID, planet, currentResearchID, planetRepo)
		if err != nil {
			return fmt.Errorf("generateEventForResearch(): %w", err)
		}

		err = planetRepo.CreateResearchEvent(ctx, researchEvent)
		if err != nil {
			return fmt.Errorf("planetStorage.CreateResearchEvent(): %w", err)
		}

		return nil
	})
}

func (s *Service) generateEventForResearch(ctx context.Context, user_id uuid.UUID, planetID uuid.UUID, currentResearchID consts.ResearchID, planetRepo TxStorages) (models.ResearchEvent, error) {
	nextLvlResearchID, err := s.registry.GetResearchNextLvlID(currentResearchID)
	if err != nil {
		return models.ResearchEvent{}, fmt.Errorf("registry.GetResearchNextLvlID(): %w", err)
	}

	nextLvlResearchStats, err := s.registry.GetResearchStatsByID(nextLvlResearchID)
	if err != nil {
		return models.ResearchEvent{}, fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
	}

	// Calculate resources
	resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return models.ResearchEvent{}, fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
	}

	if resources.Metal < nextLvlResearchStats.MetalCost ||
		resources.Crystal < nextLvlResearchStats.CrystalCost ||
		resources.Gas < nextLvlResearchStats.GasCost {
		return models.ResearchEvent{}, models.ErrNotEnoughResources
	}

	leftResources := models.Resources{
		Metal:     resources.Metal - nextLvlResearchStats.MetalCost,
		Crystal:   resources.Crystal - nextLvlResearchStats.CrystalCost,
		Gas:       resources.Gas - nextLvlResearchStats.GasCost,
		UpdatedAt: resources.UpdatedAt,
	}

	err = planetRepo.SetResources(ctx, planetID, leftResources)
	if err != nil {
		return models.ResearchEvent{}, fmt.Errorf("planetRepo.SetResources(): %w", err)
	}

	startedAt := time.Now()
	researchEvent := models.ResearchEvent{
		UserID:     user_id,
		ResearchID: currentResearchID,
		StartedAt:  startedAt,
		FinishedAt: startedAt.Add(time.Duration(nextLvlResearchStats.ResearchTimeS) * time.Second),
	}

	return researchEvent, nil
}
