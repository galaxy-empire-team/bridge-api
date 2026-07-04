package system

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *MissionStorage) GetSystemPlanets(ctx context.Context, system models.System) ([]models.PlanetInfo, error) {
	const getSystemPlanetsQuery = `
			SELECT 
				p.id, 
				p.z,
				u.login,
				p.has_moon,
				p.attacked_at,
				COALESCE(d.metal, 0),
				COALESCE(d.crystal, 0)
			FROM session_beta.planets p
			JOIN session_beta.users u ON p.user_id = u.id
			LEFT JOIN session_beta.planet_debris d ON p.id = d.planet_id
			WHERE p.x = $1 AND p.y = $2;
		`

	rows, err := s.DB.Query(ctx, getSystemPlanetsQuery, system.X, system.Y)
	if err != nil {
		return nil, fmt.Errorf("DB.Query(): %w", err)
	}
	defer rows.Close()

	var (
		systemPlanets []models.PlanetInfo
		attackedAt    *time.Time
	)
	for rows.Next() {
		var planet models.PlanetInfo
		err := rows.Scan(
			&planet.ID,
			&planet.Z,
			&planet.UserLogin,
			&planet.HasMoon,
			&attackedAt,
			&planet.Debris.Metal,
			&planet.Debris.Crystal,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		if attackedAt != nil {
			planet.AttackedAt = *attackedAt
		}

		systemPlanets = append(systemPlanets, planet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return systemPlanets, nil
}
