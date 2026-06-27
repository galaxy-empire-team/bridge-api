package models

import (
	"time"
)

type Resources struct {
	Metal     uint64
	Crystal   uint64
	Gas       uint64
	UpdatedAt time.Time
}

func (r Resources) IsEmpty() bool {
	return r == Resources{}
}
