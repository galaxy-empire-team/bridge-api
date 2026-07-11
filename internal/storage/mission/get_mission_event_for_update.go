package mission

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *MissionStorage) GetMissionForUpdate(ctx context.Context, userID uuid.UUID, id uint64) (models.CancelMission, error) {
	const getMissionForUpdateQuery = `
		SELECT
			id,
			planet_from,
			planet_to_x,
			planet_to_y,
			planet_to_z,
			is_returning,
			started_at,
			finished_at
		FROM session_beta.event_missions
		WHERE id = $1 AND user_id = $2;
	`

	var cancelMission models.CancelMission
	err := s.DB.QueryRow(ctx, getMissionForUpdateQuery, id, userID).Scan(
		&cancelMission.ID,
		&cancelMission.PlanetFrom,
		&cancelMission.PlanetTo.X,
		&cancelMission.PlanetTo.Y,
		&cancelMission.PlanetTo.Z,
		&cancelMission.IsReturning,
		&cancelMission.StartedAt,
		&cancelMission.FinishedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.CancelMission{}, models.ErrMissionNotFound
		}

		return models.CancelMission{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}
	return cancelMission, nil
}
