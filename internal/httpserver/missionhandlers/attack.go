package missionhandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/httpserver/middleware"
	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func Attack(missionService MissionService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req AttackRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		err = missionService.Attack(
			c.Request.Context(),
			userID,
			req.PlanetFrom,
			toCoordinatesModel(req.PlanetTo),
			toFleetUnits(req.FleetUnitCount),
		)
		if err != nil {
			handleAttackError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "attack mission started",
		})
	}
}

func handleAttackError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrFleetCannotBeEmpty):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Err: "fleet cannot be empty for attack mission",
		})
	case errors.Is(err, models.ErrPlanetNotFound):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "planet not found at the target coordinates",
		})
	case errors.Is(err, models.ErrNotEnoughFleetUnits):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Err: "not enough fleet units",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
