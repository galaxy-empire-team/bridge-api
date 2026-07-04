package consts

import "time"

const (
	GalaxyCount          = 1
	SystemInGalaxyCount  = 3
	PlanetsInSystemCount = 15

	ZeroResearchLevel       = 0
	MaxResearchesInProgress = 1

	ZeroBuildingLevel      = 0
	MaxBuildingsInProgress = 2

	NPCTierOnePositionZ   = 3
	NPCTierTwoPositionZ   = 7
	NPCTierThreePositionZ = 13
	NPCTierOneLogin       = "Pirates tier I"
	NPCTierTwoLogin       = "Pirates tier II"
	NPCTierThreeLogin     = "Pirates tier III"

	PlanetAttackCooldown = time.Hour * 6
	NPCAttackCooldown    = time.Minute * 5

	MistPlanetCoordinateZ = 16

	SpeedMultiplierMin        = 0.3
	SpeedMultiplierMax        = 1.0
	TimeToReachNearestPlanetS = 20
	TimeToReachNearestSystemS = 30
	TimeToReachNearestGalaxyS = 1200
)
