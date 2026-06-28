package eventhandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func CancelBuildingUpgrade(eventService EventService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req CancelBuildingUpgradeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		err = eventService.CancelBuildingUpgrade(c.Request.Context(), userID, req.PlanetID, req.BuildingID)
		if err != nil {
			handleCancelBuildingUpgradeError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func handleCancelBuildingUpgradeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrPlanetDoesNotBelongToUser):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "planet does not belong to user",
		})
	case errors.Is(err, models.ErrEventIsNotScheduled):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "event is not scheduled",
		})
	case errors.Is(err, registry.ErrNotFound):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "buildingID not found on planet",
		})
	case errors.Is(err, models.ErrBoostNotFound):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "boost for user not found",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
