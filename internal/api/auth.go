package api

import (
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"
	"github.com/gofiber/fiber/v2"
)

type authAPi struct {
	userService   domain.UserService
	factorService domain.FactorService
}

func NewAuth(app *fiber.App, userService domain.UserService, factorService domain.FactorService, authMid fiber.Handler) {
	h := authAPi{
		userService:   userService,
		factorService: factorService,
	}
	v1 := app.Group("/api/v1")
	v1.Post("user/login", h.LoginUser)
	v1.Get("token/validate-token", authMid, h.ValidateToken)
	v1.Post("user/register", h.RegisterUser)
	v1.Post("user/validate-otp", h.ValidateOTP)
	v1.Post("user/create-pin", authMid, h.CreatePIN)
}

// @Summary Login User
// @Description Login a user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.AuthReq true "Authentication Request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /user/login [post]
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

// @Summary Validate Token
// @Description Validate user token
// @Tags Authentication
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /token/validate-token [get]
func (a authAPi) ValidateToken(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-users")
	return ctx.Status(200).JSON(dto.SuccessResponse{
		Code:   200,
		Status: "OK",
		Data:   user,
	})
}

// @Summary Register User
// @Description Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.UserRegisterReq true "User Registration Request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /user/register [post]
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

// @Summary Validate OTP
// @Description Validate OTP for user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.ValidateOtpReq true "OTP Validation Request"
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Router /user/validate-otp [post]
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

// @Summary Create PIN
// @Description Create a PIN for the user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreatePinReq true "Create PIN Request"
// @Security ApiKeyAuth
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Router /user/create-pin [post]
func (a authAPi) CreatePIN(ctx *fiber.Ctx) error {
	var req dto.CreatePinReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	user := ctx.Locals("x-users").(dto.UserData)
	req.UserID = user.ID

	err := a.factorService.Insert(ctx.Context(), req)
	if err != nil {
		code := util.GetHttpStatus(err)
		return ctx.Status(code).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}
	return ctx.SendStatus(200)
}
