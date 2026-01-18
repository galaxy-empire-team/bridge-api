package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"initialservice/internal/config"
)

type HttpServer struct {
	server    *gin.Engine
	apiRouter *gin.RouterGroup
}

func New(logger *zap.Logger) *HttpServer {
	gin.SetMode(gin.ReleaseMode)

	server := gin.New()

	server.Use(func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Info("transport",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("μs", time.Since(start).Microseconds()),
		)
	})

	return &HttpServer{
		server:    server,
		apiRouter: server.Group("/api/v1"),
	}
}

func (hs *HttpServer) Start(ctx context.Context, cfg config.Server, logger *zap.Logger) error {
	srv := &http.Server{
		Addr:              cfg.Endpoint,
		Handler:           hs.server,
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()

	logger.Info("--- shutting down http server ---")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil { //nolint:contextcheck
		return fmt.Errorf("server.Shutdown(): %w", err)
	}

	return nil
}
