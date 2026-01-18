package user

import (
	"initialservice/internal/db"
)

type Repository struct {
	DB db.DB
}

func New(db db.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
