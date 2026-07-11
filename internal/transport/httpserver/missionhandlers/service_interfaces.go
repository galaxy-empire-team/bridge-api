package missionhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type MissionService interface {
	GetCurrentMissions(ctx context.Context, userID uuid.UUID) ([]models.UserMission, error)
	Colonize(ctx context.Context, mission models.MissionStart) (models.UserMission, error)
	Attack(ctx context.Context, mission models.MissionStart) (models.UserMission, error)
	Spy(ctx context.Context, mission models.MissionStart) (models.UserMission, error)
	Transport(ctx context.Context, mission models.MissionStart) (models.UserMission, error)
	Recycle(ctx context.Context, mission models.MissionStart) (models.UserMission, error)
	Mist(ctx context.Context, mission models.MissionStart) (models.UserMission, error)
	CancelMission(ctx context.Context, userID uuid.UUID, missionID uint64) (models.CancelMission, error)
}
