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

type UserResearches struct {
	UserID           uuid.UUID
	Research         []consts.ResearchID
	ResearchProgress []ResearchProgressInfo
}

type ResearchProgressInfo struct {
	UserID     uuid.UUID
	ResearchID consts.ResearchID
	StartedAt  time.Time
	FinishedAt time.Time
}
