package systemtimehandlers

import (
	"time"
)

type UTCResponse struct {
	UTC time.Time `json:"utc"`
}
