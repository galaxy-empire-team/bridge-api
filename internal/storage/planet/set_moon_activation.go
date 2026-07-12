package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) SetMoonActivation(ctx context.Context, planetID uuid.UUID, activateUntill time.Time) error {
	const setMoonActivationQuery = `
		UPDATE session_beta.planet_moons
		SET 
			active_until = $2
		WHERE planet_id = $1;
	`

	cmd, err := r.DB.Exec(ctx, setMoonActivationQuery, planetID, activateUntill)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return models.ErrMoonNotFound
	}

	return nil
}
