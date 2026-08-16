package systemtime

import (
	"context"
	"time"
)

func (s *Service) GetUTC(_ context.Context) time.Time {
	return time.Now().UTC()
}
