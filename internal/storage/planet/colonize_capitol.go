package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"initialservice/internal/models"
)

func (r *Repository) ColonizeCapitol(ctx context.Context, userID uuid.UUID, planetToColonize models.Planet) error {
	const createCapitolQuery = `
		WITH user_recorses_row_create AS (
			INSERT INTO session_beta.planet_resources (planet_id, metal, crystal, gas, updated_at)
			VALUES ($1, 0, 0, 0, NOW())
		)
		INSERT INTO session_beta.planets (
			id,
			user_id,
			x,
			y,
			z,
			has_moon,
			is_capitol,
			colonized_at
		) VALUES (
			$1, --- planet.ID
			$2, --- userID
			$3,	--- planet.X
			$4, --- planet.Y
			$5, --- planet.Z
			FALSE,
			TRUE,
			NOW()
		);
	`
	
	_, err := r.DB.Pool.Exec(
		ctx,
		createCapitolQuery,
		planetToColonize.ID,
		userID,
		planetToColonize.X,
		planetToColonize.Y,
		planetToColonize.Z,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				if pgErr.ConstraintName == "planet_have_unique_x_y_z" {
					return fmt.Errorf("DB.Pool.Exec(): %w", models.ErrPlanetCoordinatesAlreadyTaken)
				}

				return fmt.Errorf("DB.Pool.Exec(): %w", models.ErrCapitolAlreadyExists)
			}
		}

		return fmt.Errorf("DB.Pool.Exec(): %w", err)
	}

	return nil
}
