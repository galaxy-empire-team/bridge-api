package missionhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type MissionService interface {
	GetCurrentMissions(ctx context.Context, userID uuid.UUID) ([]models.UserMission, error)
	Colonize(ctx context.Context, mission models.MissionStart) error
	Attack(ctx context.Context, mission models.MissionStart) error
	Spy(ctx context.Context, mission models.MissionStart) error
	Transport(ctx context.Context, mission models.MissionStart) error
}
