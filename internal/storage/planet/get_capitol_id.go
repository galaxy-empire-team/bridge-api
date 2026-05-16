package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetCapitolID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	const getPlanetCapitolQuery = `
		SELECT id
		FROM session_beta.planets
		WHERE user_id = $1 AND is_capitol = true;
	`

	var planetID uuid.UUID
	err := r.DB.QueryRow(ctx, getPlanetCapitolQuery, userID).Scan(&planetID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.UUID{}, models.ErrCapitolNotFound
	}

	return planetID, nil
}
