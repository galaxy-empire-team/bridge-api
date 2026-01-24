package planethandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpgradeBuilding(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req CreateBuildingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		err := planetService.UpgradeBuilding(c.Request.Context(), req.PlanetID, req.BuildingType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Err: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "capitol colonized successfully",
		})
	}
}
