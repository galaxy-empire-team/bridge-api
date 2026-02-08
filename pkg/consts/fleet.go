package consts

type FleetUnitID uint8
type FleetUnitType string

const (
	FleetUnitTypeFighter     FleetUnitType = "fighter"
	FleetUnitTypeCorsair     FleetUnitType = "corsair"
	FleetUnitTypeDestroyer   FleetUnitType = "destroyer"
	FleetUnitTypeCruiser     FleetUnitType = "cruiser"
	FleetUnitTypeBattleship  FleetUnitType = "battleship"
	FleetUnitTypeAnnihilator FleetUnitType = "annihilator"
)
