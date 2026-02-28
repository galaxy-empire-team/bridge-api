package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetIDByCoordinates(ctx context.Context, coordinates models.Coordinates) (uuid.UUID, error) {
	const getIDQuery = `
		SELECT 
			id
		FROM session_beta.planets
		WHERE x = $1 AND y = $2 AND z = $3;
	`

	var id uuid.UUID
	err := r.DB.QueryRow(ctx, getIDQuery, coordinates.X, coordinates.Y, coordinates.Z).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return id, nil
}
