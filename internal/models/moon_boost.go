package models

import (
	"time"

	"github.com/google/uuid"
)

type MoonInfo struct {
	PlanetID       uuid.UUID
	HasMoon        bool
	ActivateUntill time.Time
}
