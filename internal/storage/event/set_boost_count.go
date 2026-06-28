package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *EventStorage) SetBoost(ctx context.Context, userID uuid.UUID, boost models.UserBoost) error {
	const setBoostCountQuery = `
		UPDATE session_beta.user_boosts
		SET 
			count = $3, 
			updated_at = now()
		WHERE 
			user_id = $1 
		AND 
			boost_id = $2;
	`

	_, err := r.DB.Exec(ctx, setBoostCountQuery, userID, boost.ID, boost.Count)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
