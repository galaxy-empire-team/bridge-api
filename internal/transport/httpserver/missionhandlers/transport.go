package missionhandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
)

func Transport(missionService MissionService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req TransportRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		err = missionService.Transport(
			c.Request.Context(),
			models.MissionStart{
				UserID:          userID,
				PlanetFrom:      req.PlanetFrom,
				PlanetTo:        toCoordinatesModel(req.PlanetTo),
				Cargo:           toResources(req.Cargo),
				Fleet:           toFleetUnits(req.FleetUnitCount),
				SpeedMultiplier: req.SpeedMultiplier,
			},
		)
		if err != nil {
			handleTransportError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "transport mission started",
		})
	}
}

func handleTransportError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrPlanetDoesNotBelongToUser):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Err: "the planet does not belong to the user",
		})
	case errors.Is(err, models.ErrNotEnoughFleetUnits):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Err: "not enough fleet units",
		})
	case errors.Is(err, models.ErrNotEnoughResources):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Err: "not enough resources",
		})
	case errors.Is(err, models.ErrFleetCannotBeEmpty):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Err: "fleet cannot be empty for transport mission",
		})
	case errors.Is(err, models.ErrTransportCargoExceedsFleetCapacity):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Err: "transport cargo exceeds fleet capacity",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
