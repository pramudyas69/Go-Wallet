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
	v1 := app.Group("/api/v1")
	v1.Post("user/login", h.LoginUser)
	v1.Get("token/validate-token", authMid, h.ValidateToken)
	v1.Post("user/register", h.RegisterUser)
	v1.Post("user/validate-otp", h.ValidateOTP)
}

func (a authAPi) LoginUser(ctx *fiber.Ctx) error {
	var req dto.AuthReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	token, err := a.userService.Authenticate(ctx.Context(), req)
	if err != nil {
		code := util.GetHttpStatus(err)
		return ctx.Status(code).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(dto.SuccessResponse{
		Code:   200,
		Status: "OK",
		Data:   token,
	})
}

func (a authAPi) ValidateToken(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-users")
	return ctx.Status(200).JSON(dto.SuccessResponse{
		Code:   200,
		Status: "OK",
		Data:   user,
	})
}

func (a authAPi) RegisterUser(ctx *fiber.Ctx) error {
	var req dto.UserRegisterReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	//fmt.Println(req)
	res, err := a.userService.Register(ctx.Context(), req)
	if err != nil {
		code := util.GetHttpStatus(err)
		return ctx.Status(code).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}
	return ctx.Status(200).JSON(dto.SuccessResponse{
		Code:   200,
		Status: "OK",
		Data:   res,
	})
}

func (a authAPi) ValidateOTP(ctx *fiber.Ctx) error {
	var req dto.ValidateOtpReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	err := a.userService.ValidateOTP(ctx.Context(), req)
	if err != nil {
		code := util.GetHttpStatus(err)
		return ctx.Status(code).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}
	return ctx.SendStatus(200)
}
