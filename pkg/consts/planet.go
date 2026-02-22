package consts

import "github.com/google/uuid"

type PlanetID uuid.UUID

type PlanetPositionX uint8
type PlanetPositionY uint16
type PlanetPositionZ uint8

func (p PlanetPositionX) ToUint8() uint8 {
	return uint8(p)
}

func (p PlanetPositionY) ToUint16() uint16 {
	return uint16(p)
}

func (p PlanetPositionZ) ToUint8() uint8 {
	return uint8(p)
}
