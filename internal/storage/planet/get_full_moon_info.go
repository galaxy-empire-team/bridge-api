package planet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetFullMoonInfo(ctx context.Context, planetID uuid.UUID) (models.MoonInfo, error) {
	const getMoonInfoQuery = `
		SELECT 
			p.has_moon,
			pm.active_until
		FROM session_beta.planets p
		LEFT JOIN session_beta.planet_moons pm ON p.id = pm.planet_id
		WHERE p.id = $1;
	`

	var activateUntill *time.Time
	moonInfo := models.MoonInfo{PlanetID: planetID}
	err := r.DB.QueryRow(ctx, getMoonInfoQuery, planetID).Scan(
		&moonInfo.HasMoon,
		&activateUntill,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.MoonInfo{}, nil
		}

		return models.MoonInfo{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	if activateUntill != nil {
		moonInfo.ActivateUntill = *activateUntill
	}

	return moonInfo, nil
}
