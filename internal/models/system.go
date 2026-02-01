package models

import "github.com/google/uuid"

type System struct {
	X uint64
	Y uint64
}

type SystemPlanets struct {
	System  System
	Planets []PlanetInfo
}

type PlanetInfo struct {
	ID        uuid.UUID
	Z         uint64
	Type      string
	UserLogin string
	HasMoon   bool
}
