package systemtimehandlers

import (
	"context"
	"time"
)

type SystemTimeService interface {
	GetUTC(ctx context.Context) time.Time
}
