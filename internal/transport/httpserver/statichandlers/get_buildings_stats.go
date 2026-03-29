package statichandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBuildingStats(systemService StaticService) func(c *gin.Context) {
	return func(c *gin.Context) {
		buildingStats, err := systemService.GetBuildingStats(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, toBuildingStatsResponse(buildingStats))
	}
}
