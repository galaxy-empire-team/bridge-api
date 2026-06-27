package planethandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
)

func GetUserResources(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		resources, err := planetService.GetUserResources(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Err: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, toUserResourcesResponse(resources))
	}
}
