package eventhandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
)

func BoostFleetConstruction(eventService EventService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req BoostFleetConstructionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		eventFinishTime, err := eventService.BoostFleetConstruction(
			c.Request.Context(),
			userID,
			req.PlanetID,
			models.UserBoost{
				ID:    req.Boost.ID,
				Count: req.Boost.Count,
			},
		)
		if err != nil {
			handleBoostFleetConstructionError(c, err)
			return
		}

		c.JSON(http.StatusOK, FinishTimeResponse{
			StartedAt:  eventFinishTime.StartedAt.UTC(),
			FinishedAt: eventFinishTime.FinishedAt.UTC(),
		})
	}
}

func handleBoostFleetConstructionError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrPlanetDoesNotBelongToUser):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "planet does not belong to user",
		})
	case errors.Is(err, models.ErrEventIsNotScheduled):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "event is not scheduled",
		})
	case errors.Is(err, models.ErrNotEnoughBoosts):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "not enough boosts",
		})
	case errors.Is(err, models.ErrBoostNotFound):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "boost for user not found",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
