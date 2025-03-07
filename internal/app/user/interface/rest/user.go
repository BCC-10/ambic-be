package rest

import (
	"ambic/internal/app/user/usecase"
	"ambic/internal/domain/dto"
	"ambic/internal/infra/helper"
	res "ambic/internal/infra/response"
	"ambic/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	Validator   *validator.Validate
	UserUsecase usecase.UserUsecaseItf
	helper      helper.HelperIf
}

func NewUserHandler(routerGroup fiber.Router, userUsecase usecase.UserUsecaseItf, validator *validator.Validate, m middleware.MiddlewareIf, helper helper.HelperIf) {
	UserHandler := UserHandler{
		Validator:   validator,
		UserUsecase: userUsecase,
		helper:      helper,
	}

	routerGroup = routerGroup.Group("/users")
	routerGroup.Get("/profile", m.Authentication, UserHandler.GetUserProfile)
	routerGroup.Patch("/update", m.Authentication, UserHandler.UpdateUser)
}

func (h UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	req := new(dto.UpdateUserRequest)
	if err := h.helper.FormParser(ctx, req); err != nil {
		return res.BadRequest(ctx, err.Error())
	}

	if err := h.Validator.Struct(req); err != nil {
		return res.ValidationError(ctx, nil, err)
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	if err := h.UserUsecase.UpdateUser(userId, *req); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.UpdateSuccess, nil)
}

func (h UserHandler) GetUserProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(uuid.UUID)

	user, _err := h.UserUsecase.GetUserProfile(userId)
	if _err != nil {
		return res.Error(ctx, _err)
	}

	return res.SuccessResponse(ctx, res.ShowUserSuccess, fiber.Map{
		"user": user,
	})
}
