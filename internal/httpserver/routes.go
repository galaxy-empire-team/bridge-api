package httpserver

import (
	"initialservice/internal/httpserver/planethandlers"
	"initialservice/internal/httpserver/userhandlers"
)

func (hs *HttpServer) RegisterRoutes(
	userService userhandlers.UserService, 
	planetService planethandlers.PlanetService,
) {
	hs.apiRouter.POST("/user/create", userhandlers.CreateUser(userService))

	hs.apiRouter.POST("/planet/capitol", planethandlers.GetCapitol(planetService))
	hs.apiRouter.POST("/planet/capitol/colonize", planethandlers.ColonizeCapitol(planetService))
}
