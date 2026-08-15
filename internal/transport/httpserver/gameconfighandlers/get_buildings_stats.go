package gameconfighandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBuildingStats(gameConfigService GameConfigService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var versionRequest VersionRequest
		if err := c.ShouldBindJSON(&versionRequest); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		gameConfig, err := gameConfigService.GetConfig(c.Request.Context(), versionRequest.Version)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Err: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, GameConfig{
			Version: gameConfig.Version,
			Config:  gameConfig.Config,
		})
	}
}
