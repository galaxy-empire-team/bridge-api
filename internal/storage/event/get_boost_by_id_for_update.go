package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *EventStorage) GetBoostByIDForUpdate(ctx context.Context, userID uuid.UUID, boostID consts.BoostID) (models.UserBoost, error) {
	const getBoostByIDQuery = `
		SELECT 
			boost_id,
			count
		FROM session_beta.user_boosts
		WHERE user_id = $1 AND boost_id = $2
		FOR UPDATE;
	`

	var boost models.UserBoost
	err := r.DB.QueryRow(ctx, getBoostByIDQuery, userID, boostID).Scan(&boost.ID, &boost.Count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserBoost{}, models.ErrBoostNotFound
		}

		return models.UserBoost{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return boost, nil
}
