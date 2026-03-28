package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/httpserver/middleware"
	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func StartFleetConstruction(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req StartFleetConstructionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		err = planetService.StartFleetConstruction(c.Request.Context(), userID, req.PlanetID, models.FleetUnitCount{
			ID:    req.FleetID,
			Count: req.Count,
		})
		if err != nil {
			handleStartFleetConstructionError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "fleet construction started",
		})
	}
}

func handleStartFleetConstructionError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrPlanetDoesNotBelongToUser):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "planet does not belong to user",
		})
	case errors.Is(err, models.ErrResearchInProgress):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "research already in progress",
		})
	case errors.Is(err, registry.ErrMaxLevelReached):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "research has reached max level",
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
