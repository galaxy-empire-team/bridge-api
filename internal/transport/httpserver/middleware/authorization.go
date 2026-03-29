package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const userIDKey = "userID"

// Authorization is a middleware that handles authorization.
func Authorization() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Err: "missing authorization header",
			})

			return
		}

		// TODO: add claim after auth service is implemented
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Err: "invalid authorization header format",
			})

			return
		}

		c.Set(userIDKey, parts[1])
	}
}

func RetrieveUserID(c *gin.Context) (uuid.UUID, error) {
	userIDValue, exists := c.Get(userIDKey)
	if !exists {
		return uuid.Nil, fmt.Errorf("%s not found in context", userIDKey)
	}

	userID, ok := userIDValue.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("parse %s: value is not a string", userIDKey)
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse %s: %w", userIDKey, err)
	}

	return parsedUUID, nil
}
