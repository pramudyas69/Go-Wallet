package api

import (
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type midtransApi struct {
	midtransService domain.MidtransService
	topUpService    domain.TopUpService
}

func NewMidtrans(app *fiber.App, authMid fiber.Handler, midtransService domain.MidtransService, topUpService domain.TopUpService) {
	h := midtransApi{
		midtransService: midtransService,
		topUpService:    topUpService,
	}
	v1 := app.Group("/api/v1")
	v1.Post("/midtrans/payment-callback", h.paymentHandlerNotification)
}

func (m midtransApi) paymentHandlerNotification(ctx *fiber.Ctx) error {
	var notificationPayload map[string]interface{}
	fmt.Println("Start")
	if err := ctx.BodyParser(&notificationPayload); err != nil {
		return ctx.SendStatus(400)
	}

	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		return ctx.SendStatus(400)
	}

	success, err := m.midtransService.VerifyPayment(ctx.Context(), orderId)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	if success {
		err = m.topUpService.ConfirmedTopUp(ctx.Context(), orderId)
		if err != nil {
			return ctx.Status(util.GetHttpStatus(err)).JSON(dto.ErrorResponse{
				Message: err.Error(),
			})
		}
		return ctx.SendStatus(200)
	}

	return ctx.SendStatus(400)
}
