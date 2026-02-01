package missionhandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func Colonize(missionService MissionService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req ColonizeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		err := missionService.Colonize(c.Request.Context(), req.UserID, req.PlanetFrom, toCoordinatesModel(req.PlanetTo))
		if err != nil {
			handleColonizeError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "colonization mission started",
		})
	}
}

func handleColonizeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrColonizePlanetAlreadyExists):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "planet already exists at the target coordinates",
		})
	case errors.Is(err, models.ErrPlanetDoesNotBelongToUser):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Err: "the planet does not belong to the user",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: "internal server error",
		})
	}
}
