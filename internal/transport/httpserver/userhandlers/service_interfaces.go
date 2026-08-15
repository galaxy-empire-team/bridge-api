package userhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserID(ctx context.Context, login string) (uuid.UUID, error)
}
