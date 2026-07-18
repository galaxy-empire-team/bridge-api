package rating

import (
	"github.com/galaxy-empire-team/bridge-api/internal/db"
)

type RatingStorage struct {
	DB db.DB
}

func New(db db.DB) *RatingStorage {
	return &RatingStorage{
		DB: db,
	}
}
