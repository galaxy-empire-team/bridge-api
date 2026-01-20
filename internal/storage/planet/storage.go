package planet

import (
	"initialservice/internal/db"
)

// Embed txManager requires different naming -> can't use 'storage' storage name :()
type PlanetStorage struct {
	DB db.DB
}

func New(db db.DB) *PlanetStorage {
	return &PlanetStorage{
		DB: db,
	}
}
