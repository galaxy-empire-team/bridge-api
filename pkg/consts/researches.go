package consts

type ResearchID uint16

func (r ResearchID) ToUint16() uint16 {
	return uint16(r)
}

type ResearchLevel uint8
type ResearchType string

const (
	ResearchTypeColonizeTechnology                 ResearchType = "colonize_technology"
	ResearchTypeIndustrialTechnology               ResearchType = "industrial_technology"
	ResearchTypeLogisticsTechnology                ResearchType = "logistics"
	ResearchTypeConstructionOptimizationTechnology ResearchType = "construction_optimization"

	ResearchTypeLootingTechnology ResearchType = "looting_technology"
	ResearchTypeWeaponTechnology  ResearchType = "weapon_tech"
	ResearchTypeArmorTechnology   ResearchType = "armor_tech"
	ResearchTypeSpyTechnology     ResearchType = "spy_technology"
)
