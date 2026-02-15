package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/httpserver/middleware"
	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func CreateCapitol(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		err = planetService.CreateCapitol(c.Request.Context(), userID)
		if err != nil {
			handleColonizeCapitolError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "capitol colonized successfully",
		})
	}
}

func handleColonizeCapitolError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrCapitolAlreadyExists):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "capitol planet already exists for user",
		})
	case errors.Is(err, models.ErrPlanetCoordinatesAlreadyTaken):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "planet coordinates are already taken",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
