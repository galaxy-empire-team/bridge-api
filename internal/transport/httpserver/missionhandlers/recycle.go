package missionhandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
)

func Recycle(missionService MissionService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req RecycleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		mission, err := missionService.Recycle(
			c.Request.Context(),
			models.MissionStart{
				UserID:          userID,
				PlanetFrom:      req.PlanetFrom,
				PlanetTo:        toCoordinatesModel(req.PlanetTo),
				Fleet:           toFleetUnits(req.FleetUnitCount),
				SpeedMultiplier: req.SpeedMultiplier,
			},
		)
		if err != nil {
			handleRecycleError(c, err)
			return
		}

		c.JSON(http.StatusOK, fromUserMission(mission))
	}
}

func handleRecycleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrFleetCannotBeEmpty):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Err: "fleet cannot be empty for recycle mission",
		})
	case errors.Is(err, models.ErrPlanetNotFound):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "planet not found at the target coordinates",
		})
	case errors.Is(err, models.ErrNotEnoughFleetUnits):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Err: "not enough fleet units",
		})
	case errors.Is(err, models.ErrNoDebrisFound):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Err: "no debris found on planet",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
