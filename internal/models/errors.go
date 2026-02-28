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
	ErrBuildingNotFound              = errors.New("building not found on planet")
	ErrNotEnoughResources            = errors.New("not enough resources")
	ErrEventIsAlreadyScheduled       = errors.New("event is already scheduled")
	ErrTooManyBuildingsInProgress    = errors.New("too many buildings are already in progress")
	ErrNoPlanetsFound                = errors.New("no planets found for user")
	ErrPlanetNotFound                = errors.New("planet not found")
	ErrFleetNotFound                 = errors.New("fleet not found")
	ErrBuildingAlreadyExists         = errors.New("building already exists on planet")

	// Mission errors
	ErrColonizePlanetAlreadyExists        = errors.New("planet already exists at the target coordinates")
	ErrPlanetDoesNotBelongToUser          = errors.New("the planet does not belong to the user")
	ErrFleetIDNotExists                   = errors.New("fleet unit with given ID does not exist")
	ErrNotEnoughFleetUnits                = errors.New("not enough fleet units on planet")
	ErrInvalidInput                       = errors.New("request fleet has invalid stuct or data")
	ErrInvalidShipTypeForSpyMission       = errors.New("invalid ship type for spy mission")
	ErrTransportCargoExceedsFleetCapacity = errors.New("transport cargo exceeds fleet capacity")
	ErrFleetCannotBeEmpty                 = errors.New("fleet cannot be empty for attack mission")
	ErrColonizationNotAvailable           = errors.New("colonization not available")
)
