package user

import (
	"context"

	"initialservice/internal/models"
)

type userRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

type Service struct {
	userRepo userRepository
}

func New(repo userRepository) *Service {
	return &Service{userRepo: repo}
}
