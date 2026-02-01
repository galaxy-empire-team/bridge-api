package mission

import (
	"github.com/galaxy-empire-team/bridge-api/internal/db"
)

// Embed txManager requires different naming -> can't use 'storage' storage name :()
type MissionStorage struct {
	DB db.DB
}

func New(db db.DB) *MissionStorage {
	return &MissionStorage{
		DB: db,
	}
}
