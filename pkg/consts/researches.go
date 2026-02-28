package consts

type ResearchID uint16

func (r ResearchID) ToUint16() uint16 {
	return uint16(r)
}

type ResearchLevel uint8
type ResearchType string

const (
	ResearchTypeColonizeTechnology ResearchType = "colonize_technology"

	ResearchTypeIndustrialTechnology     ResearchType = "industrial_technology"
	ResearchTypeLogistics                ResearchType = "logistics"
	ResearchTypeConstructionOptimization ResearchType = "construction_optimization"
	ResearchTypePlanetDefense            ResearchType = "planet_defense"

	ResearchTypeWeaponTech      ResearchType = "weapon_tech"
	ResearchTypeArmorTech       ResearchType = "armor_tech"
	ResearchTypeCombatProtocols ResearchType = "combat_protocols"
	ResearchTypeSpyTechnology   ResearchType = "spy_technology"
)
