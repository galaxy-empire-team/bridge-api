package userhandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func GetUserID(userService UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		userID, err := userService.GetUserID(c.Request.Context(), req.Login)
		if err != nil {
			handleGetUserIDError(c, err)
			return
		}

		c.JSON(http.StatusOK, UserResponse{
			UserID: userID.String(),
		})
	}
}

func handleGetUserIDError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrUserNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Err: "user not found",
		})
	case errors.Is(err, models.ErrInvalidLogin):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "invalid login",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: "internal server error",
		})
	}
}
