package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"initialservice/internal/models"
)

func (r *Repository) GetCapitol(ctx context.Context, userID uuid.UUID) (models.Planet, error) {
	const getCapitolQuery = `
		SELECT 
			p.id, 
			p.x, 
			p.y, 
			p.z, 
			p.has_moon, 
			p.colonized_at,
			pr.metal, 
			pr.crystal, 
			pr.gas,
			pr.updated_at
		FROM session_beta.planets p
		JOIN session_beta.planet_resources pr ON p.id = pr.planet_id
		WHERE p.is_capitol = TRUE AND p.user_id = $1;
	`

	var planet Planet

	err := r.DB.Pool.QueryRow(ctx, getCapitolQuery, userID).Scan(
		&planet.ID,
		&planet.X,
		&planet.Y,
		&planet.Z,
		&planet.HasMoon,
		&planet.ColonizedAt,
		&planet.Resources.Metal,
		&planet.Resources.Crystal,
		&planet.Resources.Gas,
		&planet.Resources.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Planet{}, fmt.Errorf("DB.Pool.QueryRow(): %w", models.ErrCapitolNotFound)
		}

		return models.Planet{}, fmt.Errorf("DB.Pool.QueryRow().Scan(): %w", err)
	}

	return toModelPlanet(planet), nil
}
