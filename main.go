package main

import (
	"e-wallet/internal/api"
	"e-wallet/internal/component"
	"e-wallet/internal/config"
	"e-wallet/internal/middleware"
	"e-wallet/internal/repository"
	"e-wallet/internal/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := component.GetDatabaseConnectiom(cnf)
	cacheConnection := component.GetCacheConnection()

	userRepository := repository.NewUser(dbConnection)

	userService := service.NewUser(userRepository, cacheConnection)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()
	api.NewAuth(app, userService, authMid)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
