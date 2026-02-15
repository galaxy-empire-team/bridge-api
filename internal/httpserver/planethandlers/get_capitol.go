package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/httpserver/middleware"
	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func GetCapitol(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		userCapitolPlanet, err := planetService.GetCapitol(c.Request.Context(), userID)
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
			Err: err.Error(),
		})
	}
}
