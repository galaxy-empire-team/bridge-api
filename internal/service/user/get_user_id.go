package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetUserID(ctx context.Context, login string) (uuid.UUID, error) {
	if login == "" {
		return uuid.Nil, models.ErrInvalidLogin
	}

	userID, err := s.userStorage.GetUserIDByLogin(ctx, login)
	if err != nil {
		return uuid.Nil, fmt.Errorf("userStorage.GetUserIDByLogin(): %w", err)
	}

	return userID, nil
}
