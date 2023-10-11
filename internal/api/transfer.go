package api

import (
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"
	"github.com/gofiber/fiber/v2"
)

type transferApi struct {
	transactionService domain.TransactionService
	factorService      domain.FactorService
}

func NewTransfer(app *fiber.App,
	authMid fiber.Handler,
	transactionService domain.TransactionService,
	factorService domain.FactorService) {
	h := transferApi{
		transactionService: transactionService,
		factorService:      factorService,
	}
	v1 := app.Group("/api/v1")
	v1.Post("transfer/inquiry", authMid, h.TransferInquiry)
	v1.Post("transfer/execute", authMid, h.TransferExecute)
}

// @Summary Transfer Inquiry
// @Description Initiate a transfer inquiry.
// @Tags Transfer
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param transferInquiryReq body dto.TransferInquiryReq true "Transfer inquiry request payload"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,401,500 {object} dto.ErrorResponse
// @Router /transfer/inquiry [post]
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

// @Summary Execute Transfer
// @Description Execute a transfer based on the inquiry.
// @Tags Transfer
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param transferExecuteReq body dto.TransferExecuteReq true "Transfer execute request payload"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,401,500 {object} dto.ErrorResponse
// @Router /transfer/execute [post]
func (t transferApi) TransferExecute(ctx *fiber.Ctx) error {
	var req dto.TransferExecuteReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	user := ctx.Locals("x-users").(dto.UserData)

	err := t.factorService.ValidatePin(ctx.Context(), dto.ValidatePinReq{
		UserID: user.ID,
		Pin:    req.PIN,
	})
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = t.transactionService.TransferExecute(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}
	return ctx.SendStatus(200)
}
