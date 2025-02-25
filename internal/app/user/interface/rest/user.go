package rest

import (
	"ambic/internal/app/user/usecase"
	"ambic/internal/domain/dto"
	"ambic/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Validator   *validator.Validate
	UserUsecase usecase.UserUsecaseItf
	Middleware  middleware.MiddlewareIf
}

func NewUserHandler(routerGroup fiber.Router, userUsecase usecase.UserUsecaseItf, validator *validator.Validate, middleware middleware.MiddlewareIf) {
	UserHandler := UserHandler{
		Validator:   validator,
		UserUsecase: userUsecase,
		Middleware:  middleware,
	}

	routerGroup = routerGroup.Group("/users")
	routerGroup.Post("/register", UserHandler.Register)
	routerGroup.Post("/login", UserHandler.Login)
	routerGroup.Post("/request-otp", UserHandler.RequestOTP)
	routerGroup.Post("verify-otp", UserHandler.VerifyOTP)
}

func (h UserHandler) Register(ctx *fiber.Ctx) error {
	user := new(dto.Register)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	if err := h.Validator.Struct(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.UserUsecase.Register(*user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.UserUsecase.RequestOTP(dto.RequestOTP{Email: user.Email}); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"payload": user,
	})
}

func (h UserHandler) RequestOTP(ctx *fiber.Ctx) error {
	requestOTP := new(dto.RequestOTP)
	if err := ctx.BodyParser(requestOTP); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	if err := h.Validator.Struct(requestOTP); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.UserUsecase.RequestOTP(*requestOTP); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OTP sent successfully",
	})
}

func (h UserHandler) VerifyOTP(ctx *fiber.Ctx) error {
	verifyOTP := new(dto.VerifyOTP)
	if err := ctx.BodyParser(verifyOTP); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	if err := h.Validator.Struct(verifyOTP); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.UserUsecase.VerifyOTP(*verifyOTP); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OTP verified successfully",
	})
}

func (h UserHandler) Login(ctx *fiber.Ctx) error {
	user := new(dto.Login)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	if err := h.Validator.Struct(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	token, err := h.UserUsecase.Login(*user)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}
