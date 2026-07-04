package missionhandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
)

func CancelMissionEvent(missionService MissionService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req CancelMissionEventRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		updatedMissionEvent, err := missionService.CancelMission(c.Request.Context(), userID, req.MissionDBID)
		if err != nil {
			handleCancelMissionEventError(c, err)
			return
		}

		c.JSON(http.StatusOK, CancelMissionEventResponse{
			ID:          updatedMissionEvent.ID,
			IsReturning: updatedMissionEvent.IsReturning,
			StartedAt:   updatedMissionEvent.StartedAt,
			FinishedAt:  updatedMissionEvent.FinishedAt,
		})
	}
}

func handleCancelMissionEventError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrMissionIsReturning):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "mission is already returning",
		})
	case errors.Is(err, models.ErrMissionIsReturning):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "mission already returning",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
