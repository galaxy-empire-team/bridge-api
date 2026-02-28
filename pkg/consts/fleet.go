package consts

type FleetUnitID uint8

func (b FleetUnitID) ToUint8() uint8 {
	return uint8(b)
}

type FleetUnitType string

const (
	FleetUnitTypeScout       FleetUnitType = "scout"
	FleetUnitTypeFighter     FleetUnitType = "fighter"
	FleetUnitTypeCorsair     FleetUnitType = "corsair"
	FleetUnitTypeDestroyer   FleetUnitType = "destroyer"
	FleetUnitTypeCruiser     FleetUnitType = "cruiser"
	FleetUnitTypeBattleship  FleetUnitType = "battleship"
	FleetUnitTypeAnnihilator FleetUnitType = "annihilator"
	FleetUnitTypeColonyShip  FleetUnitType = "colonyship"
	FleetUnitTypeRecycler    FleetUnitType = "recycler"
)
