package user

import (
	"context"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type userStorage interface {
	CreateUser(ctx context.Context, user models.User) error
}

type Service struct {
	userStorage userStorage
}

func New(userStorage userStorage) *Service {
	return &Service{userStorage: userStorage}
}
