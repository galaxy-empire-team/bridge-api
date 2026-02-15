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
	// Rewrite post > get after authorization
	hs.apiRouter.GET("/planet", planethandlers.GetPlanet(planetService))
	hs.apiRouter.GET("/planet/capitol", planethandlers.GetCapitol(planetService))
	hs.apiRouter.GET("/planet/all", planethandlers.GetAllUserPlanets(planetService))
	hs.apiRouter.GET("/planet/fleet", planethandlers.GetFleet(planetService))
	hs.apiRouter.POST("/planet/capitol/colonize", planethandlers.CreateCapitol(planetService))

	hs.apiRouter.GET("/planet/building/stats", planethandlers.GetBuildingStats(planetService))
	hs.apiRouter.POST("/planet/building/upgrade", planethandlers.UpgradeBuilding(planetService))

	// ----- Mission Routes -----
	hs.apiRouter.GET("/mission/all", missionhandlers.GetCurrentMissions(missionService))
	hs.apiRouter.POST("/mission/colonize", missionhandlers.Colonize(missionService))
	hs.apiRouter.POST("/mission/attack", missionhandlers.Attack(missionService))

	// ----- System Routes -----
	hs.apiRouter.GET("/system/planets", systemhandlers.GetSystemPlanets(systemService))
}
