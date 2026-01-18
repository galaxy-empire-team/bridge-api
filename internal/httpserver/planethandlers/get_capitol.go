package planethandlers

import (
	"net/http"
	"errors"

	"github.com/gin-gonic/gin"

	"initialservice/internal/models"
)

func GetCapitol(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		userCapitolPlanet, err := planetService.GetCapitolPlanet(c.Request.Context(), req.UserID)
		if err != nil {
			handleGetCapitolError(c, err)
			return
		}

		c.JSON(http.StatusOK, toTransportPlanet(userCapitolPlanet))
	}
}

func handleGetCapitolError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrCapitolNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "capitol planet for user not found",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
		})
	}
}
