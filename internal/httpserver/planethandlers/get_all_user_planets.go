package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func GetAllUserPlanets(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req UserIDRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		userPlanets, err := planetService.GetAllUserPlanets(c.Request.Context(), req.UserID)
		if err != nil {
			handleGetAllPlanetsError(c, err)
			return
		}

		c.JSON(http.StatusOK, toTransportPlanets(userPlanets))
	}
}

func handleGetAllPlanetsError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNoPlanetsFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Err: "no planets found for user",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
