package gameconfig

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type configStorage interface {
	GetLastConfigVersion(ctx context.Context) (models.GameConfig, error)
}

type Service struct {
	configStorage configStorage
	config        models.GameConfig
}

func New(ctx context.Context, configStorage configStorage) (*Service, error) {
	s := &Service{configStorage: configStorage}

	config, err := s.configStorage.GetLastConfigVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("configStorage.GetLastConfigVersion(): %w", err)
	}

	s.config = config

	return s, nil
}
