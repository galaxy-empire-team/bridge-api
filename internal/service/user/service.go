package user

import (
	"context"

	userstorage "initialservice/internal/storage/user"
)

type userRepository interface {
	CreateUser(ctx context.Context, user userstorage.User) (userstorage.User, error)
}

type Service struct {
	userRepo userRepository
}

func New(repo userRepository) *Service {
	return &Service{userRepo: repo}
}
