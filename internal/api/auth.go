package api

import (
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"
	"github.com/gofiber/fiber/v2"
)

type authAPi struct {
	userService domain.UserService
}

func NewAuth(app *fiber.App, userService domain.UserService, authMid fiber.Handler) {
	h := authAPi{userService: userService}

	app.Post("token/generate", h.GenerateToken)
	app.Get("token/validate", authMid, h.ValidateToken)
}

func (a authAPi) GenerateToken(ctx *fiber.Ctx) error {
	var req dto.AuthReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(util.GetHttpStatus(err))
	}

	token, err := a.userService.Authenticate(ctx.Context(), req)
	if err != nil {
		return ctx.SendStatus(util.GetHttpStatus(err))
	}

	return ctx.Status(200).JSON(token)
}

func (a authAPi) ValidateToken(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-users")
	return ctx.Status(200).JSON(user)
}
