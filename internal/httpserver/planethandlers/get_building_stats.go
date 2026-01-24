package planethandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBuildingStats(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req GetBuildStatsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: err.Error(),
			})
			return
		}

		buildingStats, err := planetService.GetBuildingStats(c.Request.Context(), req.BuildingType, req.Level)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, toTransportBuildingStats(buildingStats))
	}
}
