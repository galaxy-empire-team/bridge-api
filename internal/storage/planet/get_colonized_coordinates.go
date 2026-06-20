package planet

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *PlanetStorage) GetColonizedCoordinates(ctx context.Context, system models.Coordinates) ([]consts.PlanetPositionZ, error) {
	const getCoordinatesQuery = `
		SELECT 
			z
		FROM session_beta.planets
		WHERE x = $1 AND y = $2;
	`

	var coordinates []consts.PlanetPositionZ
	rows, err := r.DB.Query(ctx, getCoordinatesQuery, system.X, system.Y)
	if err != nil {
		return nil, fmt.Errorf("DB.Query(): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var coordinate consts.PlanetPositionZ
		err := rows.Scan(&coordinate)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}
		coordinates = append(coordinates, coordinate)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return coordinates, nil
}
