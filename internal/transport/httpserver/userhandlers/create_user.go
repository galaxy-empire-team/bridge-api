package userhandlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

const (
	minLoginLength = 3
	maxLoginLength = 50
)

type UserRequest struct {
	Login string `json:"login"`
}

type UserResponse struct {
	UserID string `json:"userID"`
}

type ErrorResponse struct {
	Err string `json:"err"`
}

func CreateUser(userService UserService) func(c *gin.Context) {
	loginChecker := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	return func(c *gin.Context) {
		var req UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: "invalid request body",
			})
			return
		}

		req.Login = strings.TrimSpace(req.Login)
		if err := validateLogin(req.Login, loginChecker); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Err: err.Error(),
			})
			return
		}

		createdUser, err := userService.CreateUser(c.Request.Context(), models.User{
			Login: req.Login,
		})
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.JSON(http.StatusCreated, UserResponse{
			UserID: createdUser.ID.String(),
		})
	}
}

func validateLogin(login string, checker *regexp.Regexp) error {
	if login == "" {
		return fmt.Errorf("login cannot be empty")
	}

	if len(login) < minLoginLength {
		return fmt.Errorf("login must be at least %d characters long", minLoginLength)
	}

	if len(login) > maxLoginLength {
		return fmt.Errorf("login must be at most %d characters long", maxLoginLength)
	}

	if !checker.MatchString(login) {
		return fmt.Errorf("login must contain only English letters, digits, '_' or '-' characters")
	}

	return nil
}

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrUserAlreadyExists):
		c.JSON(http.StatusConflict, ErrorResponse{
			Err: "user with this login already exists",
		})
	default:
		// Internal server error for unhandled cases
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Err: "internal server error",
		})
	}
}
