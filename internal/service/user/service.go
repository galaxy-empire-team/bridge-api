package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type userStorage interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserIDByLogin(ctx context.Context, login string) (uuid.UUID, error)
}

type Service struct {
	userStorage userStorage
}

func New(userStorage userStorage) *Service {
	return &Service{userStorage: userStorage}
}
