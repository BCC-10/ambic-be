package rest

import (
	"ambic/internal/app/user/usecase"
	"ambic/internal/domain/dto"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Validator   *validator.Validate
	UserUsecase usecase.UserUsecaseItf
}

func NewUserHandler(routerGroup fiber.Router, userUsecase usecase.UserUsecaseItf, validator *validator.Validate) {
	UserHandler := UserHandler{
		Validator:   validator,
		UserUsecase: userUsecase,
	}

	routerGroup = routerGroup.Group("/auth")
	routerGroup.Post("/register", UserHandler.Register)
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

	err := h.UserUsecase.Register(*user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}
