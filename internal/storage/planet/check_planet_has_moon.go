package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *PlanetStorage) CheckPlanetHasMoon(ctx context.Context, planetID uuid.UUID) (bool, error) {
	const getPlanetHasMoonQuery = `
		SELECT 
			has_moon
		FROM session_beta.planets
		WHERE id = $1;
	`

	var hasMoon bool
	err := s.DB.QueryRow(ctx, getPlanetHasMoonQuery, planetID).Scan(&hasMoon)
	if err != nil {
		return false, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return hasMoon, nil
}
