package consts

type BoostID uint8

func (b BoostID) ToUint8() uint8 {
	return uint8(b)
}

type BoostTier uint8

func (b BoostTier) ToUint8() uint8 {
	return uint8(b)
}

const (
	BoostTier1 BoostTier = 1
	BoostTier2 BoostTier = 2
	BoostTier3 BoostTier = 3

	BoostID1 BoostID = 1
	BoostID2 BoostID = 2
	BoostID3 BoostID = 3
)
