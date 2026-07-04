package systemhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type SystemService interface {
	GetSystemPlanets(ctx context.Context, userID uuid.UUID, system models.System) (models.SystemPlanets, error)
}
