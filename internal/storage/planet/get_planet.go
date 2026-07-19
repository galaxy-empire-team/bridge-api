package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetPlanet(ctx context.Context, planetID uuid.UUID) (models.Planet, error) {
	const getPlanetsQuery = `
		SELECT 
			p.id, 
			p.x,
			p.y,
			p.z,
			r.metal,
			r.crystal,
			r.gas,
			r.updated_at,
			p.has_moon,
			pm.active_until,
			p.is_capitol,
			p.colonized_at,
			p.updated_at
		FROM session_beta.planets p
		JOIN session_beta.planet_resources r ON p.id = r.planet_id
		JOIN session_beta.planet_moons pm ON p.id = pm.planet_id
		WHERE p.id = $1;
	`

	var planet models.Planet
	var colonizedAt, updatedAt time.Time
	err := r.DB.QueryRow(ctx, getPlanetsQuery, planetID).Scan(
		&planet.ID,
		&planet.Coordinates.X,
		&planet.Coordinates.Y,
		&planet.Coordinates.Z,
		&planet.Resources.Metal,
		&planet.Resources.Crystal,
		&planet.Resources.Gas,
		&planet.Resources.UpdatedAt,
		&planet.HasMoon,
		&planet.MoonActiveUntil,
		&planet.IsCapitol,
		&colonizedAt,
		&updatedAt,
	)

	planet.ColonizedAt = colonizedAt.UTC()
	planet.UpdatedAt = updatedAt.UTC()

	if err != nil {
		return models.Planet{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return planet, nil
}
