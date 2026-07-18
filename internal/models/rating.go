package models

import "github.com/google/uuid"

type RatingPlayer struct {
	UserID         uuid.UUID
	Login          string
	SpentResources uint64
	FleetPower     uint64
	Rank           uint32
	PreviousRank   uint32
}

type Rating struct {
	Top  []RatingPlayer
	Near []RatingPlayer
}

type Ratings struct {
	User  Rating
	Fleet Rating
}
