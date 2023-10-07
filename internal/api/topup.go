package api

import (
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"
	"github.com/gofiber/fiber/v2"
)

type topupApi struct {
	topupService domain.TopUpService
}

func NewTopUp(app *fiber.App, authMid fiber.Handler, topupService domain.TopUpService) {
	h := topupApi{topupService}
	v1 := app.Group("/api/v1")
	v1.Post("/topup/initialize", authMid, h.InitializeTopUp)
}

func (t topupApi) InitializeTopUp(ctx *fiber.Ctx) error {
	var req dto.TopUpReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	user := ctx.Locals("x-users").(dto.UserData)
	req.UserID = user.ID

	topUp, err := t.topupService.InitializeTopUp(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(dto.SuccessResponse{
		Code:   200,
		Status: "OK",
		Data:   topUp,
	})
}
