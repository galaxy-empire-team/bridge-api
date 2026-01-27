package models

import (
	"time"

	"github.com/google/uuid"
)

type BuildingType string

const (
	BuildingTypeMetalMine   BuildingType = "metal_mine"
	BuildingTypeCrystalMine BuildingType = "crystal_mine"
	BuildingTypeGasMine     BuildingType = "gas_mine"
)

type BuildEvent struct {
	PlanetID     uuid.UUID
	BuildingType BuildingType
	StartedAt    time.Time
	FinishedAt   time.Time
}

type BuildingStats struct {
	Level                uint8
	Type                 BuildingType
	MetalCost            uint64
	CrystalCost          uint64
	GasCost              uint64
	MetalPerSecond       uint64
	CrystalPerSecond     uint64
	GasPerSecond         uint64
	Bonuses              *string
	UpgradeTimeInSeconds uint16
}

type BuildingInfo struct {
	Level            uint8
	Type             BuildingType
	MetalPerSecond   uint64
	CrystalPerSecond uint64
	GasPerSecond     uint64
	Bonuses          *string
	UpdatedAt        time.Time
	FinishedAt       time.Time
}

func GetMines() []BuildingType {
	return []BuildingType{
		BuildingTypeMetalMine,
		BuildingTypeCrystalMine,
		BuildingTypeGasMine,
	}
}

func GetAllBuildings() []BuildingType {
	var buildings []BuildingType

	buildings = append(buildings, GetMines()...)

	return buildings
}

func IsValidBuildingType(buildingType BuildingType) bool {
	for _, bt := range GetAllBuildings() {
		if bt == buildingType {
			return true
		}
	}

	return false
}
