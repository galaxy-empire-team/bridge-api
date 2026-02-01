package models

import "errors"

var (
	// User errors
	ErrUserAlreadyExists = errors.New("user already exists")

	// Planet errors
	ErrCapitolAlreadyExists          = errors.New("capitol planet already exists")
	ErrCapitolNotFound               = errors.New("capitol planet not found")
	ErrPlanetCoordinatesAlreadyTaken = errors.New("planet coordinates already taken")
	ErrBuildingMaxLevelReached       = errors.New("building has reached max level")
	ErrBuildingInvalidLevel          = errors.New("invalid building level")
	ErrBuildTypeInvalid              = errors.New("invalid building type")
	ErrNotEnoughResources            = errors.New("not enough resources to upgrade building")
	ErrEventIsAlreadyScheduled       = errors.New("event is already scheduled")
	ErrTooManyBuildingsInProgress    = errors.New("too many buildings are already in progress")

	// Mission errors
	ErrColonizePlanetAlreadyExists = errors.New("planet already exists at the target coordinates")
	ErrPlanetDoesNotBelongToUser   = errors.New("the planet does not belong to the user")
)
