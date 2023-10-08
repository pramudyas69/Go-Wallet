package main

import (
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/api"
	"e-wallet/internal/component"
	"e-wallet/internal/config"
	"e-wallet/internal/middleware"
	"e-wallet/internal/repository"
	"e-wallet/internal/service"
	"e-wallet/internal/sse"
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
	hub := &dto.Hub{
		NotificationChannel: make(map[int64]chan dto.NotificationData),
	}

	userRepository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)
	notificationRepository := repository.NewNotification(dbConnection)
	templateRepository := repository.NewTemplate(dbConnection)
	topUpRepository := repository.NewTopUp(dbConnection)
	factorRepository := repository.NewFactor(dbConnection)

	emailService := service.NewEmail(cnf)
	userService := service.NewUser(userRepository, cacheConnection, emailService, accountRepository, utilInterface, jwtInterface)
	notificationService := service.NewNotification(notificationRepository, templateRepository, hub)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection, emailService, userRepository, utilInterface, notificationService)
	midtransService := service.NewMidtransService(cnf)
	topUpService := service.NewTopUpService(topUpRepository, midtransService, accountRepository, notificationService, transactionRepository)
	factorService := service.NewFactor(factorRepository)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()
	api.NewAuth(app, userService, factorService, authMid)
	api.NewTransfer(app, authMid, transactionService, factorService)
	api.NewNotification(app, authMid, notificationService)
	sse.NewNotification(app, authMid, hub)
	api.NewTopUp(app, authMid, topUpService)
	api.NewMidtrans(app, authMid, midtransService, topUpService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
