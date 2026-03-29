package missionhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type MissionService interface {
	GetCurrentMissions(ctx context.Context, userID uuid.UUID) ([]models.UserMission, error)
	Colonize(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates, cargo models.Resources, fleet []models.FleetUnitCount) error
	Attack(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates, fleet []models.FleetUnitCount) error
	Spy(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates, fleet []models.FleetUnitCount) error
	Transport(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates, cargo models.Resources, fleet []models.FleetUnitCount) error
}
