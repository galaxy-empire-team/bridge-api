package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/httpserver/middleware"
	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func GetPlanet(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req PlanetIDRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		planet, err := planetService.GetPlanet(c.Request.Context(), userID, req.PlanetID)
		if err != nil {
			handleGetPlanetError(c, err)
			return
		}

		c.JSON(http.StatusOK, toPlanetResponse(planet))
	}
}

func handleGetPlanetError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrPlanetNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Err: "planet not found",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
