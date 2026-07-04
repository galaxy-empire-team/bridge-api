package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *MissionStorage) CheckNPCMissionExists(ctx context.Context, userID uuid.UUID, z consts.PlanetPositionZ) (bool, error) {
	const checkNPCMissionExistsQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM session_beta.event_missions
			WHERE user_id = $1 AND planet_to_z = $2 AND is_returning = false
		);
	`

	var exists bool
	err := s.DB.QueryRow(ctx, checkNPCMissionExistsQuery, userID, z).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return exists, nil
}
