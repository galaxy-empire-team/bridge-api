package main

import (
	"fmt"
	"log"

	"github.com/galaxy-empire-team/bridge-api/internal/app"
	"github.com/galaxy-empire-team/bridge-api/internal/config"
	"github.com/galaxy-empire-team/bridge-api/internal/db"
	"github.com/galaxy-empire-team/bridge-api/internal/httpserver"
	missionservice "github.com/galaxy-empire-team/bridge-api/internal/service/mission"
	planetservice "github.com/galaxy-empire-team/bridge-api/internal/service/planet"
	systemservice "github.com/galaxy-empire-team/bridge-api/internal/service/system"
	userservice "github.com/galaxy-empire-team/bridge-api/internal/service/user"
	missionstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/mission"
	planetstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/planet"
	systemstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/system"
	"github.com/galaxy-empire-team/bridge-api/internal/storage/txmanager"
	userstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/user"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
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
	systemStorage := systemstorage.New(db)
	missionStorage := missionstorage.New(db)

	// initialize registry
	reg, err := registry.New(ctx, db.Pool)
	if err != nil {
		return fmt.Errorf("registry.New(): %w", err)
	}

	// initialize services
	userService := userservice.New(userStorage)
	planetService := planetservice.New(txManager, planetStorage, reg)
	missionService := missionservice.New(txManager, planetStorage, missionStorage, reg)
	systemService := systemservice.New(systemStorage)

	// initialize http server
	httpServer := httpserver.New(app.ComponentLogger("httpserver"))

	httpServer.RegisterRoutes(userService, planetService, missionService, systemService)

	err = httpServer.Start(ctx, cfg.Server)
	if err != nil {
		return fmt.Errorf("httpServer.Start(): %w", err)
	}

	return nil
}
