package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/galaxy-empire-team/bridge-api/internal/config"

	"go.uber.org/zap"
)

const (
	shutdownTimeout = 5 * time.Second
)

type App struct {
	cancelFn      context.CancelFunc
	gracefulFuncs []func(context.Context) error

	logger *zap.Logger
}

func New(cfg config.App) (context.Context, *App, error) {
	ctx, cancelFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger, err := newLogger(cfg)
	if err != nil {
		return context.Background(), nil, fmt.Errorf("newLogger(): %w", err)
	}

	go func() {
		<-ctx.Done()
		logger.Sync() // nolint:errcheck, gosec
	}()

	return ctx, &App{
		cancelFn: cancelFn,
		logger:   logger,
	}, nil
}

func (a *App) ComponentLogger(component string) *zap.Logger {
	return a.logger.With(zap.String("component", component))
}

func (a *App) AddGracefulFunc(fn func(context.Context) error) {
	a.gracefulFuncs = append(a.gracefulFuncs, fn)
}

func (a *App) WaitAndShutdown(ctx context.Context) error {
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	for _, fn := range a.gracefulFuncs {
		if err := fn(shutdownCtx); err != nil {
			a.logger.Error("graceful shutdown func failed", zap.Error(err))
		}
	}

	a.logger.Info("app shutdown completed")

	return nil
}
