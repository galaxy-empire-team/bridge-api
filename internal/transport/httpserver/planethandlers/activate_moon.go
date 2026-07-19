package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
)

func ActivateMoon(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req ActivateMoonRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		moonInfo, err := planetService.ActivateMoon(
			c.Request.Context(),
			userID,
			req.PlanetID,
			req.BoostID,
			req.Count,
		)
		if err != nil {
			handleActivateMoonError(c, err)
			return
		}

		c.JSON(http.StatusOK, PlanetMoonInfoResponse{
			PlanetID:       moonInfo.PlanetID,
			HasMoon:        moonInfo.HasMoon,
			ActivateUntill: moonInfo.ActivateUntill.UTC(),
		})
	}
}

func handleActivateMoonError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrPlanetDoesNotBelongToUser):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "planet does not belong to user",
		})
	case errors.Is(err, models.ErrMoonNotFound):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "moon not found",
		})
	case errors.Is(err, models.ErrNotEnoughMatter):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "not enough matter",
		})
	case errors.Is(err, models.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "invalid request",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
