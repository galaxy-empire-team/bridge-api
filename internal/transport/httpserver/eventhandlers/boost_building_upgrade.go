package eventhandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
)

func BoostBuildingUpgrade(eventService EventService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req BoostBuildingUpgradeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		eventFinishTime, err := eventService.BoostBuildingUpgrade(
			c.Request.Context(),
			userID,
			req.PlanetID,
			req.BuildingID,
			models.UserBoost{
				ID:    req.Boost.ID,
				Count: req.Boost.Count,
			},
		)
		if err != nil {
			handleBoostBuildingUpgradeError(c, err)
			return
		}

		c.JSON(http.StatusOK, FinishTimeResponse{
			StartedAt:  eventFinishTime.StartedAt.UTC(),
			FinishedAt: eventFinishTime.FinishedAt.UTC(),
		})
	}
}

func handleBoostBuildingUpgradeError(c *gin.Context, err error) {
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
