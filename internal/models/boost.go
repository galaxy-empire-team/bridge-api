package models

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type UserBoost struct {
	ID    consts.BoostID
	Count uint64
}
