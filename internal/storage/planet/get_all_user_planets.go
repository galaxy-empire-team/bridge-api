package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error) {
	const getAllUserPlanetsQuery = `
		SELECT 
			id, 
			x,
			y,
			z,
			has_moon,
			is_capitol,
			colonized_at
		FROM session_beta.planets
		WHERE user_id = $1;
	`

	rows, err := r.DB.Query(ctx, getAllUserPlanetsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}
	defer rows.Close()

	var planets []models.Planet
	var colonizedAt time.Time
	for rows.Next() {
		var planet models.Planet
		err = rows.Scan(
			&planet.ID,
			&planet.Coordinates.X,
			&planet.Coordinates.Y,
			&planet.Coordinates.Z,
			&planet.HasMoon,
			&planet.IsCapitol,
			&colonizedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		planet.ColonizedAt = colonizedAt.UTC()

		planets = append(planets, planet)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return planets, nil
}
