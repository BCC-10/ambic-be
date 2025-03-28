package rest

import (
	"ambic/internal/app/auth/usecase"
	"ambic/internal/domain/dto"
	"ambic/internal/infra/limiter"
	res "ambic/internal/infra/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Validator   *validator.Validate
	AuthUsecase usecase.AuthUsecaseItf
}

func NewAuthHandler(routerGroup fiber.Router, userUsecase usecase.AuthUsecaseItf, validator *validator.Validate, limiter limiter.LimiterIf) {
	AuthHandler := AuthHandler{
		Validator:   validator,
		AuthUsecase: userUsecase,
	}

	routerGroup = routerGroup.Group("/auth")
	routerGroup.Post("/register", AuthHandler.Register)
	routerGroup.Post("/login", AuthHandler.Login)
	routerGroup.Get("/verification", limiter.Set(3, "15m"), AuthHandler.ResendVerification)
	routerGroup.Post("/verification", AuthHandler.VerifyUser)
	routerGroup.Post("/forgot-password", limiter.Set(3, "15m"), AuthHandler.ForgotPassword)
	routerGroup.Patch("/reset-password", AuthHandler.ResetPassword)
	routerGroup.Get("/google", AuthHandler.GoogleLogin)
	routerGroup.Post("/google", AuthHandler.GoogleCallback)
}

func (h AuthHandler) Register(ctx *fiber.Ctx) error {
	user := new(dto.RegisterRequest)
	if err := ctx.BodyParser(user); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(user); err != nil {
		return res.ValidationError(ctx, err)
	}

	if err := h.AuthUsecase.Register(*user); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.RegisterSuccess, nil)
}

func (h AuthHandler) ResendVerification(ctx *fiber.Ctx) error {
	requestToken := new(dto.EmailVerificationRequest)
	if err := ctx.QueryParser(requestToken); err != nil {
		return res.BadRequest(ctx, err.Error())
	}

	if err := h.Validator.Struct(requestToken); err != nil {
		return res.ValidationError(ctx, err)
	}

	if err := h.AuthUsecase.SendVerification(*requestToken); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.VerificationLinkSent, nil)
}

func (h AuthHandler) VerifyUser(ctx *fiber.Ctx) error {
	data := new(dto.VerifyUserRequest)
	if err := ctx.BodyParser(data); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return res.ValidationError(ctx, err)
	}

	if err := h.AuthUsecase.VerifyUser(*data); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.VerifySuccess, nil)
}

func (h AuthHandler) Login(ctx *fiber.Ctx) error {
	user := new(dto.LoginRequest)
	if err := ctx.BodyParser(user); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(user); err != nil {
		return res.ValidationError(ctx, err)
	}

	token, err := h.AuthUsecase.Login(*user)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.LoginSuccess, fiber.Map{
		"token": token,
	})
}

func (h AuthHandler) ForgotPassword(ctx *fiber.Ctx) error {
	data := new(dto.ForgotPasswordRequest)
	if err := ctx.BodyParser(data); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return res.ValidationError(ctx, err)
	}

	if err := h.AuthUsecase.ForgotPassword(*data); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.ForgotPasswordSuccess, nil)
}

func (h AuthHandler) ResetPassword(ctx *fiber.Ctx) error {
	data := new(dto.ResetPasswordRequest)
	if err := ctx.BodyParser(data); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return res.ValidationError(ctx, err)
	}

	if err := h.AuthUsecase.ResetPassword(*data); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.ResetPasswordSuccess, nil)
}

func (h AuthHandler) GoogleLogin(ctx *fiber.Ctx) error {
	url, err := h.AuthUsecase.GoogleLogin()
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.OAuthLoginSuccess, fiber.Map{
		"url": url,
	})
}

func (h AuthHandler) GoogleCallback(ctx *fiber.Ctx) error {
	data := new(dto.GoogleCallbackRequest)
	if err := ctx.BodyParser(data); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return res.ValidationError(ctx, err)
	}

	token, err := h.AuthUsecase.GoogleCallback(*data)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.LoginSuccess, fiber.Map{
		"token": token,
	})
}
