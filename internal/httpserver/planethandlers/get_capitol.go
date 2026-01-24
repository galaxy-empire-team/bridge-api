package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"initialservice/internal/models"
)

func GetCapitol(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req UserIDRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: err.Error(),
			})
			return
		}

		userCapitolPlanet, err := planetService.GetCapitol(c.Request.Context(), req.UserID)
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
			Err: "capitol planet for user not found",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: "internal server error" + err.Error(),
		})
	}
}
