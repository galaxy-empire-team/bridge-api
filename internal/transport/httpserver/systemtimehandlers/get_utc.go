package systemtimehandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUTC(systemTimeService SystemTimeService) func(c *gin.Context) {
	return func(c *gin.Context) {
		utc := systemTimeService.GetUTC(c.Request.Context())

		c.JSON(http.StatusOK, UTCResponse{
			UTC: utc,
		})
	}
}
