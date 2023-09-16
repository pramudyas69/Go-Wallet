package middleware

import (
	"e-wallet/domain"
	"e-wallet/internal/util"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Authenticate(userService domain.UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := strings.ReplaceAll(ctx.Get("Authorization"), "Bearer ", "")
		if token == "" {
			return ctx.SendStatus(401)
		}
		user, err := userService.ValidateToken(ctx.Context(), token)
		if err != nil {
			return ctx.SendStatus(util.GetHttpStatus(err))
		}
		ctx.Locals("x-users", user)
		return ctx.Next()
	}
}
