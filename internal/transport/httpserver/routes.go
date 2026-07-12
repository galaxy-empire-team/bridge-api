package httpserver

import (
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/eventhandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/missionhandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/planethandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/statichandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/systemhandlers"
	"github.com/galaxy-empire-team/bridge-api/internal/transport/httpserver/userhandlers"
)

func (hs *HttpServer) RegisterRoutes(
	userService userhandlers.UserService,
	planetService planethandlers.PlanetService,
	eventService eventhandlers.EventService,
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
	hs.apiRouter.GET("/planet/moon", planethandlers.GetMoonInfo(planetService))
	hs.apiRouter.POST("/planet/moon/activate", planethandlers.ActivateMoon(planetService))
	hs.apiRouter.POST("/planet/capitol/colonize", planethandlers.ColonizeCapitol(planetService))

	// ----- Event Routes -----
	hs.apiRouter.POST("/event/building/upgrade", eventhandlers.StartBuildingUpgrade(eventService))
	hs.apiRouter.POST("/event/building/cancel", eventhandlers.CancelBuildingUpgrade(eventService))
	hs.apiRouter.POST("/event/building/boost", eventhandlers.BoostBuildingUpgrade(eventService))
	hs.apiRouter.POST("/event/research/start", eventhandlers.StartResearch(eventService))
	hs.apiRouter.POST("/event/research/cancel", eventhandlers.CancelResearch(eventService))
	hs.apiRouter.POST("/event/research/boost", eventhandlers.BoostResearch(eventService))
	hs.apiRouter.POST("/event/fleet/construct", eventhandlers.StartFleetConstruction(eventService))
	hs.apiRouter.POST("/event/fleet/cancel", eventhandlers.CancelFleetConstruction(eventService))
	hs.apiRouter.POST("/event/fleet/boost", eventhandlers.BoostFleetConstruction(eventService))

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
