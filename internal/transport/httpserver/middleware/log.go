package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggingMiddleware(logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		w, ok := c.Writer.(*bodyWriter)
		if !ok {
			logger.Error("cast response writer")
			return
		}

		status := c.Writer.Status()
		if status >= http.StatusInternalServerError {
			logger.Error("processing error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Int("status", status),
				zap.Int64("μs", time.Since(start).Microseconds()),
				zap.String("body", w.buf.String()),
			)

			return
		}

		logger.Info("processing success",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("μs", time.Since(start).Microseconds()),
		)
	}
}
