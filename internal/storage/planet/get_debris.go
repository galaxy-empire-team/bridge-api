package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetDebris(ctx context.Context, coordinates models.Coordinates) (models.Resources, error) {
	const getDebrisQuery = `
		SELECT 
			metal,
			crystal
		FROM session_beta.planet_debris pd
		JOIN session_beta.planets p ON pd.planet_id = p.id
		WHERE p.x = $1 AND p.y = $2 AND p.z = $3;
	`

	var debris models.Resources
	err := r.DB.QueryRow(ctx, getDebrisQuery, coordinates.X, coordinates.Y, coordinates.Z).Scan(
		&debris.Metal,
		&debris.Crystal,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Resources{}, nil
		}

		return models.Resources{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return debris, nil
}
