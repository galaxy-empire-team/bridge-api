package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func GetFleet(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req PlanetIDRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		fleet, err := planetService.GetFleet(c.Request.Context(), req.PlanetID)
		if err != nil {
			handleGetFleetError(c, err)
			return
		}

		c.JSON(http.StatusOK, toTransportFleet(fleet))
	}
}

func handleGetFleetError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrPlanetNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Err: "fleet not found",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
