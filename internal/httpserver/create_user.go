package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"initialservice/internal/models"
)

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
	Error string `json:"error"`
}

func createUser(userService userService) func(c *gin.Context) {
	loginChecker := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	return func(c *gin.Context) {
		var req UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "invalid request body",
			})
			return
		}

		req.Login = strings.TrimSpace(req.Login)
		if err := validateLogin(req.Login, loginChecker); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
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
		return fmt.Errorf("login have to contain only english letters, digits, _ or - characters")
	}

	return nil
}

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrUserAlreadyExists):
		c.JSON(http.StatusConflict, ErrorResponse{
			Error: "user with this login already exists",
		})
	default:
		// Internal server error for unhandled cases
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
		})
	}
}
