package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UpdatePlanetResources recalculates the resources of a planet based on the time since the last update.
// Recalcs using the provided updatedTime if not nil, otherwise uses time.Now().UTC(). Use this before any operation that changes resources.
// This method is used by internal gRPC calls.
func (s *Service) UpdatePlanetResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, updatedTime *time.Time) error {
	if updatedTime == nil {
		err := s.resourceRepository.RecalcResources(ctx, userID, planetID)
		if err != nil {
			return fmt.Errorf("recalcResources(): %w", err)
		}

		return nil
	}

	err := s.resourceRepository.RecalcResourcesWithUpdatedTime(ctx, userID, planetID, *updatedTime)
	if err != nil {
		return fmt.Errorf("recalcResourcesWithUpdatedTime(): %w", err)
	}

	return nil
}
