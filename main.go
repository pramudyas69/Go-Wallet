package main

import (
	"e-wallet/domain"
	"e-wallet/internal/api"
	"e-wallet/internal/component"
	"e-wallet/internal/config"
	"e-wallet/internal/middleware"
	"e-wallet/internal/repository"
	"e-wallet/internal/service"
	"e-wallet/internal/util"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()

	dbConnection := component.GetDatabaseConnectiom(cnf)
	defer dbConnection.Close()
	//cacheConnection := component.GetCacheConnection()
	cacheConnection := repository.NewRedisClient(cnf)
	utilInterface := domain.NewUtil()
	jwtInterface := util.NewJwt(cnf)

	userRepository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)
	notificationRepository := repository.NewNotification(dbConnection)

	emailService := service.NewEmail(cnf)
	userService := service.NewUser(userRepository, cacheConnection, emailService, accountRepository, utilInterface, jwtInterface)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection, emailService, userRepository, utilInterface, notificationRepository)
	notificationService := service.NewNotification(notificationRepository)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()
	api.NewAuth(app, userService, authMid)
	api.NewTransfer(app, authMid, transactionService)
	api.NewNotification(app, authMid, notificationService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
