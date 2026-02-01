package httpserver

import (
	"github.com/galaxy-empire-team/bridge-api/internal/httpserver/missionhandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/httpserver/planethandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/httpserver/systemhandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/httpserver/userhandlers"
)

func (hs *HttpServer) RegisterRoutes(
	userService userhandlers.UserService,
	planetService planethandlers.PlanetService,
	missionService missionhandlers.MissionService,
	systemService systemhandlers.SystemService,
) {
	// ----- User Routes -----
	hs.apiRouter.POST("/user/create", userhandlers.CreateUser(userService))

	// ----- Planet Routes -----
	hs.apiRouter.POST("/planet/capitol", planethandlers.GetCapitol(planetService))
	hs.apiRouter.POST("/planet/capitol/colonize", planethandlers.CreateCapitol(planetService))

	hs.apiRouter.POST("/planet/building/upgrade", planethandlers.UpgradeBuilding(planetService))
	hs.apiRouter.POST("/planet/building/stats", planethandlers.GetBuildingStats(planetService))

	// ----- Mission Routes -----
	hs.apiRouter.POST("/mission/colonize", missionhandlers.Colonize(missionService))

	// ----- System Routes -----
	hs.apiRouter.POST("/system/planets", systemhandlers.GetSystemPlanets(systemService))
}
