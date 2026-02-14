package consts

type FleetUnitID uint8

func (b FleetUnitID) ToUint8() uint8 {
	return uint8(b)
}

type FleetUnitType string

const (
	FleetUnitTypeFighter     FleetUnitType = "fighter"
	FleetUnitTypeCorsair     FleetUnitType = "corsair"
	FleetUnitTypeDestroyer   FleetUnitType = "destroyer"
	FleetUnitTypeCruiser     FleetUnitType = "cruiser"
	FleetUnitTypeBattleship  FleetUnitType = "battleship"
	FleetUnitTypeAnnihilator FleetUnitType = "annihilator"
)
