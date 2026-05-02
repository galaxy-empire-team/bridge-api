package planet

import (
	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type CreatePlanetRequest struct {
	OperationID uint64
	Coordinates models.Coordinates
	Resources   models.Resources
	IsCapitol   bool
}
