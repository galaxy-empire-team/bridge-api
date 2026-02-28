package research

import (
	"github.com/galaxy-empire-team/bridge-api/internal/db"
)

// Embed txManager requires different naming -> can't use 'storage' storage name :()
type ResearchStorage struct {
	DB db.DB
}

func New(db db.DB) *ResearchStorage {
	return &ResearchStorage{
		DB: db,
	}
}
