package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func GetBuildingStats(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req GetBuildStatsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		buildingStats, err := planetService.GetBuildingStats(c.Request.Context(), req.BuildingType, req.Level)
		if err != nil {
			handleGetBuildingStatsError(c, err)
			return
		}

		c.JSON(http.StatusOK, toTransportBuildingStats(buildingStats))
	}
}

func handleGetBuildingStatsError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrBuildingInvalidLevel):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "invalid building level",
		})
	case errors.Is(err, models.ErrBuildTypeInvalid):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "invalid building type",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: "internal server error",
		})
	}
}
