package main

import (
	"fmt"
	"log"

	"github.com/galaxy-empire-team/bridge-api/internal/application"
	"github.com/galaxy-empire-team/bridge-api/internal/config"
	"github.com/galaxy-empire-team/bridge-api/internal/db"
	repository "github.com/galaxy-empire-team/bridge-api/internal/repository"
	eventservice "github.com/galaxy-empire-team/bridge-api/internal/service/event"
	missionservice "github.com/galaxy-empire-team/bridge-api/internal/service/mission"
	planetservice "github.com/galaxy-empire-team/bridge-api/internal/service/planet"
	staticservice "github.com/galaxy-empire-team/bridge-api/internal/service/static"
	systemservice "github.com/galaxy-empire-team/bridge-api/internal/service/system"
	userservice "github.com/galaxy-empire-team/bridge-api/internal/service/user"
	eventstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/event"
	missionstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/mission"
	planetstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/planet"
	researchstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/research"
	systemstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/system"
	"github.com/galaxy-empire-team/bridge-api/internal/storage/txmanager"
	userstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/user"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/grpcserver"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver"
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

	ctx, app, err := application.New(cfg.App)
	if err != nil {
		return fmt.Errorf("application.New(): %w", err)
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
	eventStorage := eventstorage.New(db)
	researchStorage := researchstorage.New(db)

	// initialize registry
	reg, err := registry.New(ctx, db.Pool)
	if err != nil {
		return fmt.Errorf("registry.New(): %w", err)
	}

	// initialize repositories
	resourceRepo := repository.New(txManager, planetStorage, researchStorage, reg)

	// initialize services
	userService := userservice.New(userStorage)
	planetService := planetservice.New(planetStorage, researchStorage, resourceRepo, txManager, reg, app.ComponentLogger("planetService"))
	eventService := eventservice.New(txManager, eventStorage, planetStorage, researchStorage, resourceRepo, reg)
	missionService := missionservice.New(txManager, planetStorage, missionStorage, researchStorage, resourceRepo, reg)
	systemService := systemservice.New(planetStorage, systemStorage)
	staticService := staticservice.New(reg)

	// initialize http server
	httpServer := httpserver.New(app.ComponentLogger("httpserver"))

	httpServer.RegisterRoutes(userService, planetService, eventService, missionService, systemService, staticService)

	shutdownFunc, err := httpServer.Start(ctx, cfg.HTTPServer)
	if err != nil {
		return fmt.Errorf("httpServer.Start(): %w", err)
	}
	app.AddGracefulFunc(shutdownFunc)

	// initialize gRPC server
	grpcServer := grpcserver.New(planetService, app.ComponentLogger("grpcserver"))

	shutdownFunc, err = grpcServer.Start(cfg.GRPCServer)
	if err != nil {
		return fmt.Errorf("grpcServer.Start(): %w", err)
	}
	app.AddGracefulFunc(shutdownFunc)

	return app.WaitAndShutdown(ctx)
}
