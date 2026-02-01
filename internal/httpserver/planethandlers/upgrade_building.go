package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func UpgradeBuilding(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req UpgradeBuildingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		err := planetService.UpgradeBuilding(c.Request.Context(), req.PlanetID, req.BuildingType)
		if err != nil {
			handleUpgradeBuildingError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "build started",
		})
	}
}

func handleUpgradeBuildingError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrBuildTypeInvalid):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "invalid building type",
		})
	case errors.Is(err, models.ErrTooManyBuildingsInProgress):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "too many buildings are already in progress",
		})
	case errors.Is(err, models.ErrBuildingMaxLevelReached):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "building has reached max level",
		})
	case errors.Is(err, models.ErrNotEnoughResources):
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Err: "not enough resources to upgrade building",
		})
	case errors.Is(err, models.ErrEventIsAlreadyScheduled):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "event is already scheduled",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: "internal server error",
		})
	}
}
