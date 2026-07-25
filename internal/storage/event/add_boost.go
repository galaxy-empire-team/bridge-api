package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *EventStorage) AddBoost(ctx context.Context, userID uuid.UUID, boost models.UserBoost) error {
	const addBoostQuery = `
		INSERT INTO session_beta.user_boosts (
			user_id,
			boost_id,
			count,
			updated_at
		) VALUES (
			$1,
			$2,
			$3,
			now()
		)
		ON CONFLICT (user_id, boost_id) DO UPDATE SET
			count = session_beta.user_boosts.count + excluded.count,
			updated_at = now();
	`

	_, err := r.DB.Exec(ctx, addBoostQuery, userID, boost.ID, boost.Count)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
