package models

import (
	"time"
)

type EventFinishTime struct {
	EventID    uint64
	StartedAt  time.Time
	FinishedAt time.Time
}
