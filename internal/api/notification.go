package api

import (
	"e-wallet/domain"
	"e-wallet/dto"
	"github.com/gofiber/fiber/v2"
)

type notificationApi struct {
	notificationService domain.NotificationService
}

func NewNotification(app *fiber.App, authMid fiber.Handler, notificationService domain.NotificationService) {
	h := notificationApi{
		notificationService: notificationService,
	}
	v1 := app.Group("/api/v1")
	v1.Get("notification", authMid, h.FindByUser)
}

// @Summary Get Notifications for User
// @Description Get notifications for the authenticated user.
// @Tags Notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,401,500 {object} dto.ErrorResponse
// @Router /notification [get]
func (a notificationApi) FindByUser(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-users")

	notifications, err := a.notificationService.FindByUser(ctx.Context(), user.(dto.UserData).ID)
	if err != nil {
		return ctx.Status(400).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(dto.SuccessResponse{
		Code:   200,
		Status: "OK",
		Data:   notifications,
	})
}
