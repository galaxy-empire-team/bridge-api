package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetMoonActivationForUpdate(ctx context.Context, planetID uuid.UUID) (models.MoonInfo, error) {
	const getMoonInfoQuery = `
		SELECT 
			active_until
		FROM session_beta.planet_moons
		WHERE planet_id = $1
		FOR UPDATE;
	`

	var moonInfo models.MoonInfo
	err := r.DB.QueryRow(ctx, getMoonInfoQuery, planetID).Scan(
		&moonInfo.ActivateUntill,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.MoonInfo{}, nil
		}

		return models.MoonInfo{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return moonInfo, nil
}
