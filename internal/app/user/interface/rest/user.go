package rest

import (
	"ambic/internal/app/user/usecase"
	"ambic/internal/domain/dto"
	"ambic/internal/infra/limiter"
	res "ambic/internal/infra/response"
	"ambic/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	Validator   *validator.Validate
	UserUsecase usecase.UserUsecaseItf
}

func NewUserHandler(routerGroup fiber.Router, userUsecase usecase.UserUsecaseItf, validator *validator.Validate, middleware middleware.MiddlewareIf, limiter limiter.LimiterIf) {
	UserHandler := UserHandler{
		Validator:   validator,
		UserUsecase: userUsecase,
	}

	routerGroup = routerGroup.Group("/users")
	routerGroup.Post("/register", UserHandler.Register)
	routerGroup.Post("/login", UserHandler.Login)
	routerGroup.Post("/request-otp", limiter.Set(3, "15m"), UserHandler.RequestOTP)
	routerGroup.Post("/verify", UserHandler.VerifyUser)
	routerGroup.Post("/forgot-password", limiter.Set(3, "15m"), UserHandler.ForgotPassword)
	routerGroup.Patch("/reset-password", UserHandler.ResetPassword)
	routerGroup.Patch("/:id/update", middleware.Authentication, middleware.EnsureVerified, UserHandler.UpdateUser)
}

func (h UserHandler) Register(ctx *fiber.Ctx) error {
	user := new(dto.Register)
	if err := ctx.BodyParser(user); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(user); err != nil {
		return res.BadRequest(ctx, err.Error())
	}

	if err := h.UserUsecase.Register(*user); err != nil {
		return res.Error(ctx, err)
	}

	if err := h.UserUsecase.RequestOTP(dto.RequestOTP{Email: user.Email}); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.RegisterSuccess, user.AsResponse())
}

func (h UserHandler) RequestOTP(ctx *fiber.Ctx) error {
	requestOTP := new(dto.RequestOTP)
	if err := ctx.BodyParser(requestOTP); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(requestOTP); err != nil {
		return res.BadRequest(ctx, err.Error())
	}

	if err := h.UserUsecase.RequestOTP(*requestOTP); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.OTPSent, nil)
}

func (h UserHandler) VerifyUser(ctx *fiber.Ctx) error {
	data := new(dto.VerifyOTP)
	if err := ctx.BodyParser(data); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return res.BadRequest(ctx, err.Error())
	}

	if err := h.UserUsecase.VerifyUser(*data); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.VerifySuccess, nil)
}

func (h UserHandler) Login(ctx *fiber.Ctx) error {
	user := new(dto.Login)
	if err := ctx.BodyParser(user); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(user); err != nil {
		return res.BadRequest(ctx, err.Error())
	}

	token, err := h.UserUsecase.Login(*user)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.LoginSuccess, fiber.Map{
		"token": token,
	})
}

func (h UserHandler) ForgotPassword(ctx *fiber.Ctx) error {
	data := new(dto.ForgotPassword)
	if err := ctx.BodyParser(data); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return res.BadRequest(ctx, err.Error())
	}

	if err := h.UserUsecase.ForgotPassword(*data); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.ForgotPasswordSuccess, nil)
}

func (h UserHandler) ResetPassword(ctx *fiber.Ctx) error {
	data := new(dto.ResetPassword)
	if err := ctx.BodyParser(data); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return res.BadRequest(ctx, err.Error())
	}

	if err := h.UserUsecase.ResetPassword(*data); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.ResetPasswordSuccess, nil)
}

func (h UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	user := new(dto.UpdateUser)
	if err := ctx.BodyParser(user); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(user); err != nil {
		return res.BadRequest(ctx, err.Error())
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	if err := h.UserUsecase.UpdateUser(userId, *user); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.UpdateSuccess, nil)
}
