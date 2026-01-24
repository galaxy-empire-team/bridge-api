package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

func (r *PlanetStorage) GetLocation(ctx context.Context, planetID uuid.UUID) (models.Location, error) {
	const getLocationQuery = `
		SELECT 
			x,
			y,
			z
		FROM session_beta.planets
		WHERE id = $1;
	`

	var location models.Location
	err := r.DB.QueryRow(ctx, getLocationQuery, planetID).Scan(
		&location.X,
		&location.Y,
		&location.Z,
	)
	if err != nil {
		return models.Location{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return location, nil
}
