package systemhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func GetSystemPlanets(systemService SystemService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req SystemPlanetsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		systemPlanets, err := systemService.GetSystemPlanets(c.Request.Context(), models.System{
			X: req.X,
			Y: req.Y,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, toSystemPlanetsResponse(systemPlanets))
	}
}
