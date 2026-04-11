package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
)

func StartBuildingUpgrade(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req StartBuildingUpgradeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		finishTime, err := planetService.StartBuildingUpgrade(c.Request.Context(), userID, req.PlanetID, req.BuildingID)
		if err != nil {
			handleUpgradeBuildingError(c, err)
			return
		}

		c.JSON(http.StatusOK, FinishTimeResponse{
			StartedAt:  finishTime.StartedAt.UTC(),
			FinishedAt: finishTime.FinishedAt.UTC(),
		})
	}
}

func handleUpgradeBuildingError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrTooManyBuildingsInProgress):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "too many buildings in progress",
		})
	case errors.Is(err, models.ErrBuildingNotFound):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "building not found",
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
			Err: err.Error(),
		})
	}
}
