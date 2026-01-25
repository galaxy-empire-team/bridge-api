package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return models.User{}, fmt.Errorf("uuid.NewV7(): %w", err)
	}

	userToCreate := models.User{
		ID:    id,
		Login: user.Login,
	}

	createdUser, err := s.userRepo.CreateUser(ctx, userToCreate)
	if err != nil {
		return models.User{}, fmt.Errorf("userRepo.CreateUser(): %w", err)
	}

	return createdUser, nil
}
