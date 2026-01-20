package models

import (
	"time"

	"github.com/google/uuid"
)

type BuildType string

const (
	BuildingTypeMetalMine   BuildType = "metal_mine"
	BuildingTypeCrystalMine BuildType = "crystal_mine"
	BuildingTypeGasMine     BuildType = "gas_mine"
)

type Planet struct {
	ID          uuid.UUID
	X           uint8
	Y           uint8
	Z           uint8
	Resources   Resources
	HasMoon     bool
	ColonizedAt time.Time
}

type Resources struct {
	Metal     uint64
	Crystal   uint64
	Gas       uint64
	UpdatedAt time.Time
}

type BuildStats struct {
	Level                uint8
	Type                 BuildType
	MetalCost            uint64
	CrystalCost          uint64
	GasCost              uint64
	MetalPerSecond       uint64
	CrystalPerSecond     uint64
	GasPerSecond         uint64
	Bonuses              *string
	UpgradeTimeInSeconds uint16
}

type PlanetBuildInfo struct {
	Level            uint8
	Type             BuildType
	MetalPerSecond   uint64
	CrystalPerSecond uint64
	GasPerSecond     uint64
	Bonuses          *string
	UpdatedAt        time.Time
}
