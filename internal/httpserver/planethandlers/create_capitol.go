package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func CreateCapitol(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req UserIDRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		err := planetService.CreateCapitol(c.Request.Context(), req.UserID)
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
