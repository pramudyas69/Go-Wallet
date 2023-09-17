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
	defer dbConnection.Close()
	cacheConnection := component.GetCacheConnection()

	userRepository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)

	emailService := service.NewEmail(cnf)
	userService := service.NewUser(userRepository, cacheConnection, emailService)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()
	api.NewAuth(app, userService, authMid)
	api.NewTransfer(app, authMid, transactionService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
