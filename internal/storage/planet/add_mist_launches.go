package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *PlanetStorage) AddMistLaunches(ctx context.Context, userID uuid.UUID, count uint64) error {
	const addMistLaunchesQuery = `
		INSERT INTO session_beta.user_mist_launches (
			user_id,
			count_left,
			updated_at
		) VALUES ($1, $2, now())
		ON CONFLICT (user_id) DO UPDATE SET
			count_left = user_mist_launches.count_left + excluded.count_left,
			updated_at = now();
	`

	_, err := r.DB.Exec(ctx, addMistLaunchesQuery, userID, count)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
