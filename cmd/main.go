package main

import (
	"fmt"
	"log"

	"initialservice/internal/app"
	"initialservice/internal/config"
	"initialservice/internal/db"
	"initialservice/internal/httpserver"
	planetservice "initialservice/internal/service/planet"
	userservice "initialservice/internal/service/user"
	planetstorage "initialservice/internal/storage/planet"
	"initialservice/internal/storage/txmanager"
	userstorage "initialservice/internal/storage/user"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config.New(): %w", err)
	}

	ctx, app, err := app.New(cfg.App)
	if err != nil {
		return fmt.Errorf("app.New(): %w", err)
	}

	// initialize pgx infra
	db, err := db.New(ctx, cfg.PgConn)
	if err != nil {
		return fmt.Errorf("db.New(): %w", err)
	}
	defer db.Close()

	// initialize storages
	txManager := txmanager.New(db)
	userStorage := userstorage.New(db)
	planetStorage := planetstorage.New(db)

	// initialize services
	userService := userservice.New(userStorage)
	planetService := planetservice.New(txManager, planetStorage)

	// initialize http server
	httpServer := httpserver.New(app.ComponentLogger("httpserver"))

	httpServer.RegisterRoutes(userService, planetService)

	err = httpServer.Start(ctx, cfg.Server)
	if err != nil {
		return fmt.Errorf("httpServer.Start(): %w", err)
	}

	return nil
}
