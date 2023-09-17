package api

import (
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"
	"github.com/gofiber/fiber/v2"
)

type transferApi struct {
	transactionService domain.TransactionService
}

func NewTransfer(app *fiber.App,
	authMid fiber.Handler,
	transactionService domain.TransactionService) {
	h := transferApi{transactionService: transactionService}

	app.Post("transfer/inquiry", authMid, h.TransferInquiry)
	app.Post("transfer/execute", authMid, h.TransferExecute)
}

func (t transferApi) TransferInquiry(ctx *fiber.Ctx) error {
	var req dto.TransferInquiryReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	inquiry, err := t.transactionService.TransferInquiry(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}
	return ctx.Status(200).JSON(dto.SuccessResponse{
		Code:   200,
		Status: "OK",
		Data:   inquiry,
	})
}

func (t transferApi) TransferExecute(ctx *fiber.Ctx) error {
	var req dto.TransferExecuteReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	err := t.transactionService.TransferExecute(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}
	return ctx.SendStatus(200)
}
