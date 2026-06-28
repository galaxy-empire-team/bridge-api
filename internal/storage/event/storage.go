package event

import (
	"github.com/galaxy-empire-team/bridge-api/internal/db"
)

type EventStorage struct {
	DB db.DB
}

func New(db db.DB) *EventStorage {
	return &EventStorage{
		DB: db,
	}
}
