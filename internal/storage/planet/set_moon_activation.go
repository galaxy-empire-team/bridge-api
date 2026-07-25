package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *PlanetStorage) SetMoonActivation(ctx context.Context, planetID uuid.UUID, activateUntill time.Time) error {
	const setMoonActivationQuery = `
		INSERT INTO session_beta.planet_moons (planet_id, active_until)
		VALUES ($1, $2)
		ON CONFLICT (planet_id) DO UPDATE SET
			active_until = excluded.active_until;
	`

	_, err := r.DB.Exec(ctx, setMoonActivationQuery, planetID, activateUntill)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
