package sse

import (
	"bufio"
	"e-wallet/dto"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type notificationSse struct {
	hub *dto.Hub
}

func NewNotification(app *fiber.App, authMid fiber.Handler, hub *dto.Hub) {
	h := notificationSse{
		hub: hub,
	}
	v1 := app.Group("/api/v1")
	v1.Get("sse/notification-stream", authMid, h.StreamNotification)
}

func (s notificationSse) StreamNotification(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/event-stream")

	user := ctx.Locals("x-users").(dto.UserData)
	s.hub.NotificationChannel[user.ID] = make(chan dto.NotificationData)

	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		event := fmt.Sprintf("event: %s\n"+
			"data: \n\n", "initial")
		_, _ = fmt.Fprint(w, event)
		_ = w.Flush()

		for notification := range s.hub.NotificationChannel[user.ID] {
			data, _ := json.Marshal(notification)

			event = fmt.Sprintf("event: %s\n"+
				"data: %s\n\n", "notification-updated", data)

			_, _ = fmt.Fprint(w, event)
			_ = w.Flush()
		}
	})
	return nil
}
