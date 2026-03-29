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

	"github.com/galaxy-empire-team/bridge-api/internal/config"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/middleware"
)

type HttpServer struct {
	server    *gin.Engine
	apiRouter *gin.RouterGroup

	logger *zap.Logger
}

func New(logger *zap.Logger) *HttpServer {
	gin.SetMode(gin.ReleaseMode)

	server := gin.New()

	server.Use(
		gin.Recovery(),
		middleware.Authorization(),
		middleware.UseCustomWriter(),
		middleware.HideInternalError(logger),
		middleware.LoggingMiddleware(logger),
	)

	return &HttpServer{
		server:    server,
		apiRouter: server.Group("/api/v1"),
		logger:    logger,
	}
}

func (s *HttpServer) Start(ctx context.Context, cfg config.HTTPServer) (func(context.Context) error, error) {
	srv := &http.Server{
		Addr:              cfg.Endpoint,
		Handler:           s.server,
		ReadTimeout:       120 * time.Second,
		ReadHeaderTimeout: 120 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		s.logger.Info("---  http server started  ---", zap.String("endpoint", cfg.Endpoint))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	return func(ctx context.Context) error {
		if err := srv.Shutdown(ctx); err != nil { //nolint:contextcheck
			return fmt.Errorf("server.Shutdown(): %w", err)
		}

		return nil
	}, nil
}
