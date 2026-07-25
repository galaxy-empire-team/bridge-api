package traderhandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func BuyItem(traderService TraderService) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := middleware.RetrieveUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Err: err.Error(),
			})

			return
		}

		var req BuyItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		err = traderService.Buy(c.Request.Context(), userID, req.PlanetID, req.ItemID)
		if err != nil {
			handleBuyItemError(c, err)
			return
		}

		c.JSON(http.StatusOK, BuyItemResponse{
			Message: "success",
		})
	}
}

func handleBuyItemError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNotEnoughResources):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "not enough resources",
		})
	case errors.Is(err, models.ErrNotEnoughDoreye):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "not enough doreye",
		})
	case errors.Is(err, models.ErrMoonAlreadyExists):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "moon already exists",
		})
	case errors.Is(err, registry.ErrNotFound):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "trader item not found",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: err.Error(),
		})
	}
}
