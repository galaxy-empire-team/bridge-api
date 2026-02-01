package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetCoordinates(ctx context.Context, planetID uuid.UUID) (models.Coordinates, error) {
	const getCoordinatesQuery = `
		SELECT 
			x,
			y,
			z
		FROM session_beta.planets
		WHERE id = $1;
	`

	var coordinates models.Coordinates
	err := r.DB.QueryRow(ctx, getCoordinatesQuery, planetID).Scan(
		&coordinates.X,
		&coordinates.Y,
		&coordinates.Z,
	)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return coordinates, nil
}
