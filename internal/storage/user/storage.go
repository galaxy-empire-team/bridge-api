package user

import (
	"github.com/galaxy-empire-team/bridge-api/internal/db"
)

// Embed txManager requires different naming -> can't use 'storage' storage name :()
type UserStorage struct {
	DB db.DB
}

func New(db db.DB) *UserStorage {
	return &UserStorage{
		DB: db,
	}
}
