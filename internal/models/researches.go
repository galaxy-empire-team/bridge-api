package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type ResearchEvent struct {
	UserID     uuid.UUID
	ResearchID consts.ResearchID
	StartedAt  time.Time
	FinishedAt time.Time
}
