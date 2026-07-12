package consts

type MoonBoostID uint8

func (b MoonBoostID) ToUint8() uint8 {
	return uint8(b)
}

const (
	MoonBoostID1 MoonBoostID = 1
)
