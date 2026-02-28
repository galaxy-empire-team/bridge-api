package planethandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/httpserver/middleware"
)

func GetFleet(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req PlanetIDRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		fleet, err := planetService.GetFleet(c.Request.Context(), userID, req.PlanetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Err: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, toFleetResponse(fleet))
	}
}
