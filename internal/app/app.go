package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"initialservice/internal/config"

	"go.uber.org/zap"
)

type App struct {
	cancelFn context.CancelFunc

	logger *zap.Logger
}

func New(cfg config.App) (context.Context, *App, error) {
	ctx, cancelFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger, err := newLogger(cfg.Logger)
	if err != nil {
		return context.Background(), nil, fmt.Errorf("newLogger(): %w", err)
	}

	go func() {
		<-ctx.Done()
		logger.Sync() // nolint:errcheck
	}()

	return ctx, &App{
		cancelFn: cancelFn,
		logger:   logger,
	}, nil
}

func (a *App) ComponentLogger(component string) *zap.Logger {
	return a.logger.With(zap.String("component", component))
}
