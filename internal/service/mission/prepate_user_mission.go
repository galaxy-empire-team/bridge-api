package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) prepareUserMission(ctx context.Context, mission models.MissionStart, missionType consts.MissionType) (models.UserMission, error) {
	missionID, err := s.registry.GetMissionIDByType(missionType)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("registry.GetMissionIDByType(): %w", err)
	}

	planetFromCoordinates, err := s.planetStorage.GetCoordinates(ctx, mission.PlanetFrom)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("planetStorage.GetCoordinates(): %w", err)
	}

	missionDuration, err := s.calculateMissionDuration(durationInput{
		PlanetFrom:      planetFromCoordinates,
		PlanetTo:        mission.PlanetTo,
		Fleet:           mission.Fleet,
		SpeedMultiplier: mission.SpeedMultiplier,
		Type:            missionType,
	})
	if err != nil {
		return models.UserMission{}, fmt.Errorf("calculateMissionDuration(): %w", err)
	}

	startedAt := time.Now().UTC()
	finishedAt := startedAt.Add(missionDuration)

	userMission := models.UserMission{
		MissionID:   missionID,
		PlanetFrom:  planetFromCoordinates,
		PlanetTo:    mission.PlanetTo,
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
		IsReturning: false,
	}

	return userMission, nil
}
