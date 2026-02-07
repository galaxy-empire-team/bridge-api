package registry

import (
	"errors"
)

var (
	ErrInvalidBuildingType  = errors.New("invalid building type")
	ErrInvalidBuildingLevel = errors.New("invalid building level")
	ErrNotFound             = errors.New("not found in registry")
	ErrMaxLevelReached      = errors.New("building has reached max level")
)
