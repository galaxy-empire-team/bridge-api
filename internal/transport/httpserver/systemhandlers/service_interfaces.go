package systemhandlers

import (
	"context"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type SystemService interface {
	GetSystemPlanets(ctx context.Context, system models.System) (models.SystemPlanets, error)
}
