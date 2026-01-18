package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"initialservice/internal/config"
	"initialservice/internal/db"
	"initialservice/internal/httpserver"
	"initialservice/internal/logger"
	userservice "initialservice/internal/service/user"
	userstorage "initialservice/internal/storage/user"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config.New(): %w", err)
	}

	logger, err := logger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("logger.New(): %w", err)
	}
	defer logger.Sync() //nolint:errcheck

	db, err := db.New(ctx, cfg.PgConn)
	if err != nil {
		return fmt.Errorf("db.New(): %w", err)
	}
	defer db.Pool.Close()

	// initialize storages
	userStorage := userstorage.New(db)

	// initialize services
	userService := userservice.New(userStorage)

	// initialize http server
	httpServer := httpserver.New(logger)

	httpServer.RegisterUserRoutes(userService)

	logger.Info("--- Start app ---")

	err = httpServer.Start(ctx, cfg.Server, logger)
	if err != nil {
		return fmt.Errorf("httpServer.Start(): %w", err)
	}

	return nil
}
