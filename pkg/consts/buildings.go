package consts

type BuildingID uint8

func (b BuildingID) ToUint8() uint8 {
	return uint8(b)
}

type BuildingLevel uint8
type BuildingType string

const (
	BuildingTypeMetalMine     BuildingType = "metal_mine"
	BuildingTypeCrystalMine   BuildingType = "crystal_mine"
	BuildingTypeGasMine       BuildingType = "gas_mine"
	BuildingTypeSpaceport     BuildingType = "space_port"
	BuildingTypeResearchLab   BuildingType = "laboratory"
	BuildingTypeDefenseCenter BuildingType = "robot_factory"
)

func GetMineTypes() []BuildingType {
	return []BuildingType{
		BuildingTypeMetalMine,
		BuildingTypeCrystalMine,
		BuildingTypeGasMine,
	}
}

func GetBuildingTypes() []BuildingType {
	buildings := []BuildingType{
		BuildingTypeSpaceport,
		BuildingTypeResearchLab,
		BuildingTypeDefenseCenter,
	}

	buildings = append(buildings, GetMineTypes()...)

	return buildings
}

func IsValidBuildingType(buildingType BuildingType) bool {
	for _, bt := range GetBuildingTypes() {
		if bt == buildingType {
			return true
		}
	}

	return false
}
