package gameconfig

import (
	"context"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetConfig(ctx context.Context, version uint16) (models.GameConfig, error) {
	if s.config.Version != version {
		return s.config, nil
	}

	return models.GameConfig{
		Version: version,
	}, nil
}
