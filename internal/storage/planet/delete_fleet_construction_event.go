package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *PlanetStorage) DeleteFleetConstructionEvent(ctx context.Context, planetID uuid.UUID) (models.Resources, error) {
	const deleteFleetConstructionEventQuery = `
			DELETE FROM session_beta.event_fleet_constructions
			WHERE planet_id = $1
			RETURNING metal_cost, crystal_cost, gas_cost;
		`

	var resources models.Resources
	err := s.DB.QueryRow(ctx, deleteFleetConstructionEventQuery, planetID).Scan(&resources.Metal, &resources.Crystal, &resources.Gas)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Resources{}, models.ErrEventIsNotScheduled
		}

		return models.Resources{}, fmt.Errorf("DB.QueryRow().Scan(): %w", err)
	}

	return resources, nil
}
