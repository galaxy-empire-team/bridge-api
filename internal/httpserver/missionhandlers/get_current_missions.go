package missionhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCurrentMissions(missionService MissionService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req UserIDRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		missions, err := missionService.GetCurrentMissions(c.Request.Context(), req.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Err: "internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, fromUserMissions(missions))
	}
}
