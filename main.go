package main

import (
	_ "e-wallet/docs"
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
	"github.com/gofiber/swagger"
)

// @title E-Wallet Open API
// @version 2.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9090
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))

	api.NewAuth(app, userService, factorService, authMid)
	api.NewTransfer(app, authMid, transactionService, factorService)
	api.NewNotification(app, authMid, notificationService)
	sse.NewNotification(app, authMid, hub)
	api.NewTopUp(app, authMid, topUpService)
	api.NewMidtrans(app, authMid, midtransService, topUpService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
