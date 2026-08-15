package gameconfig

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Storage) GetLastConfigVersion(ctx context.Context) (models.GameConfig, error) {
	const getConfigByVersionQuery = `
		SELECT 
			version, config 
		FROM session_beta.h_configs 
		ORDER BY version 
		DESC LIMIT 1;
	`

	var gameConfig models.GameConfig
	err := s.DB.QueryRow(ctx, getConfigByVersionQuery).Scan(&gameConfig.Version, &gameConfig.Config)
	if err != nil {
		return models.GameConfig{}, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return gameConfig, nil
}
