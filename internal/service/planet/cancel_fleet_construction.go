package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Service) CancelFleetConstruction(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error {
	if err := s.repository.CheckPlanetOwner(ctx, userID, planetID); err != nil {
		return fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		resources, err := planetRepo.DeleteFleetConstructionEvent(ctx, planetID)
		if err != nil {
			return fmt.Errorf("planetStorage.DeleteBuildingEvent(): %w", err)
		}

		err = planetRepo.AddResources(ctx, planetID, resources)
		if err != nil {
			return fmt.Errorf("planetStorage.AddResources(): %w", err)
		}

		return nil
	})
}
