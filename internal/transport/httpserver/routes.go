package httpserver

import (
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/missionhandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/planethandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/statichandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/systemhandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/userhandlers"
)

func (hs *HttpServer) RegisterRoutes(
	userService userhandlers.UserService,
	planetService planethandlers.PlanetService,
	missionService missionhandlers.MissionService,
	systemService systemhandlers.SystemService,
	staticsService statichandlers.StaticService,
) {
	// ----- User Routes -----
	hs.apiRouter.POST("/user/create", userhandlers.CreateUser(userService))

	// ----- Planet Routes -----
	hs.apiRouter.GET("/planet", planethandlers.GetPlanet(planetService))
	hs.apiRouter.GET("/planet/capitol", planethandlers.GetCapitol(planetService))
	hs.apiRouter.GET("/planet/all", planethandlers.GetAllUserPlanets(planetService))
	hs.apiRouter.GET("/planet/fleet", planethandlers.GetFleet(planetService))
	hs.apiRouter.GET("/planet/researches", planethandlers.GetResearches(planetService))
	hs.apiRouter.GET("/planet/resources", planethandlers.GetPlanetResources(planetService))
	hs.apiRouter.GET("/planet/userresources", planethandlers.GetUserResources(planetService))
	hs.apiRouter.GET("/planet/buildings", planethandlers.GetPlanetBuildings(planetService))
	hs.apiRouter.POST("/planet/capitol/colonize", planethandlers.ColonizeCapitol(planetService))
	hs.apiRouter.POST("/planet/building/upgrade", planethandlers.StartBuildingUpgrade(planetService))
	hs.apiRouter.POST("/planet/building/cancel", planethandlers.CancelBuildingUpgrade(planetService))
	hs.apiRouter.POST("/planet/research/start", planethandlers.StartResearch(planetService))
	hs.apiRouter.POST("/planet/research/cancel", planethandlers.CancelResearch(planetService))
	hs.apiRouter.POST("/planet/fleet/construct", planethandlers.StartFleetConstruction(planetService))
	hs.apiRouter.POST("/planet/fleet/cancel", planethandlers.CancelFleetConstruction(planetService))

	// ----- Mission Routes -----
	hs.apiRouter.GET("/mission/all", missionhandlers.GetCurrentMissions(missionService))
	hs.apiRouter.POST("/mission/colonize", missionhandlers.Colonize(missionService))
	hs.apiRouter.POST("/mission/attack", missionhandlers.Attack(missionService))
	hs.apiRouter.POST("/mission/spy", missionhandlers.Spy(missionService))
	hs.apiRouter.POST("/mission/transport", missionhandlers.Transport(missionService))
	hs.apiRouter.POST("/mission/recycle", missionhandlers.Recycle(missionService))
	hs.apiRouter.POST("/mission/mist", missionhandlers.Mist(missionService))
	hs.apiRouter.POST("/mission/cancel", missionhandlers.CancelMissionEvent(missionService))

	// ----- System Routes -----
	hs.apiRouter.GET("/system/planets", systemhandlers.GetSystemPlanets(systemService))

	// ----- Static Routes -----
	hs.apiRouter.GET("/static/buildings", statichandlers.GetBuildingStats(staticsService))
}
