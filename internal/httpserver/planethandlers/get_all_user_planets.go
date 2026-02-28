package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/httpserver/middleware"
	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func GetAllUserPlanets(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		userPlanets, err := planetService.GetAllUserPlanets(c.Request.Context(), userID)
		if err != nil {
			handleGetAllPlanetsError(c, err)
			return
		}

		c.JSON(http.StatusOK, toUserPlanetsResponse(userPlanets))
	}
}

func handleGetAllPlanetsError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNoPlanetsFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Err: "no planets found",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
