package missionhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type MissionService interface {
	Colonize(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates) error
	Attack(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates, fleet []models.PlanetFleetUnitCount) error
	GetCurrentMissions(ctx context.Context, userID uuid.UUID) ([]models.UserMission, error)
}
