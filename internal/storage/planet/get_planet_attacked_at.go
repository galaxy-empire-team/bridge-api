package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetPlanetAttackedAt(ctx context.Context, coordinate models.Coordinates) (*time.Time, error) {
	const getPlanetAttackedAtQuery = `
		SELECT 
			attacked_at
		FROM session_beta.planets
		WHERE x = $1 AND y = $2 AND z = $3;
	`

	var time *time.Time

	err := r.DB.QueryRow(ctx, getPlanetAttackedAtQuery, coordinate.X, coordinate.Y, coordinate.Z).Scan(&time)
	if err != nil {
		return nil, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return time, nil
}
