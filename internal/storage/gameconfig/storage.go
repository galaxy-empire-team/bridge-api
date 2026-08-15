package gameconfig

import (
	"github.com/galaxy-empire-team/bridge-api/internal/db"
)

type Storage struct {
	DB db.DB
}

func New(db db.DB) *Storage {
	return &Storage{
		DB: db,
	}
}
