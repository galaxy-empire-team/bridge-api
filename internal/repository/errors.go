package repository

import "errors"

var (
	ErrInvalidResearchType = errors.New("invalid research type")
	ErrInvalidBuildingType = errors.New("invalid building type")
)
