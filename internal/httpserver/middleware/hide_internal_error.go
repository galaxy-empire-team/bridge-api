package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Err string `json:"err"`
}

// HideInternalError is a middleware that hides internal server errors.
// Middleware is used after logging middleware that logs internal errors.
func HideInternalError(logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {

		c.Next()

		w, ok := c.Writer.(*bodyWriter)
		if !ok {
			logger.Error("cast response writer")
			return
		}

		status := c.Writer.Status()
		if status >= http.StatusInternalServerError {
			c.Writer = w.ResponseWriter

			c.JSON(status, ErrorResponse{
				Err: "internal server error",
			})

			return
		}

		w.ResponseWriter.WriteHeader(status)
		w.ResponseWriter.Write(w.buf.Bytes())
	}
}
