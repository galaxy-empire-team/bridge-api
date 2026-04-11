package planethandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func StartResearch(planetService PlanetService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req StartResearchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		researchProgress, err := planetService.StartResearch(c.Request.Context(), userID, req.PlanetID, req.ResearchID)
		if err != nil {
			handleStartResearchError(c, err)
			return
		}

		c.JSON(http.StatusOK, ResearchProgressResponse{
			ResearchID: researchProgress.ResearchID,
			StartedAt:  researchProgress.StartedAt,
			FinishedAt: researchProgress.FinishedAt,
		})
	}
}

func handleStartResearchError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrPlanetDoesNotBelongToUser):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "planet does not belong to user",
		})
	case errors.Is(err, models.ErrResearchInProgress):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "research already in progress",
		})
	case errors.Is(err, registry.ErrMaxLevelReached):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "research has reached max level",
		})
	case errors.Is(err, models.ErrNotEnoughResources):
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Err: "not enough resources to upgrade building",
		})
	case errors.Is(err, models.ErrEventIsAlreadyScheduled):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "event is already scheduled",
		})
	case errors.Is(err, models.ErrUserHasNotResearch):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "user has not research id to upgrade",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
