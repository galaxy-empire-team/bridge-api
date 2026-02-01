package system

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *MissionStorage) GetSystemPlanets(ctx context.Context, system models.System) (models.SystemPlanets, error) {
	const getSystemPlanetsQuery = `
			SELECT 
				p.id, 
				p.z,
				u.login,
				p.has_moon
			FROM session_beta.planets p
			JOIN session_beta.users u ON p.user_id = u.id
			WHERE p.x = $1 AND p.y = $2;
		`

	rows, err := s.DB.Query(ctx, getSystemPlanetsQuery, system.X, system.Y)
	if err != nil {
		return models.SystemPlanets{}, fmt.Errorf("DB.Query(): %w", err)
	}
	defer rows.Close()

	var systemPlanets []models.PlanetInfo
	for rows.Next() {
		var planet models.PlanetInfo
		err := rows.Scan(
			&planet.ID,
			&planet.Z,
			&planet.UserLogin,
			&planet.HasMoon,
		)
		if err != nil {
			return models.SystemPlanets{}, fmt.Errorf("rows.Scan(): %w", err)
		}

		systemPlanets = append(systemPlanets, planet)
	}

	if err := rows.Err(); err != nil {
		return models.SystemPlanets{}, fmt.Errorf("rows.Err(): %w", err)
	}

	return models.SystemPlanets{
		System:  system,
		Planets: systemPlanets,
	}, nil
}
