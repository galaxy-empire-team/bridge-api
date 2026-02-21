package consts

type MissionID uint8

func (b MissionID) ToUint8() uint8 {
	return uint8(b)
}

type MissionType string

const (
	MissionTypeColonize MissionType = "colonize"
	MissionTypeAttack   MissionType = "attack"
)
