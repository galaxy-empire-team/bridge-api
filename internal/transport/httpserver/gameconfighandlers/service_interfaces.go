package gameconfighandlers

import (
	"context"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type GameConfigService interface {
	GetConfig(ctx context.Context, version uint16) (models.GameConfig, error)
}
