package rest

import (
	"ambic/internal/app/user/usecase"
	"ambic/internal/domain/dto"
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

func NewUserHandler(routerGroup fiber.Router, userUsecase usecase.UserUsecaseItf, validator *validator.Validate, middleware middleware.MiddlewareIf) {
	UserHandler := UserHandler{
		Validator:   validator,
		UserUsecase: userUsecase,
	}

	routerGroup = routerGroup.Group("/users")
	routerGroup.Patch("/update", middleware.Authentication, middleware.EnsureVerified, UserHandler.UpdateUser)
	routerGroup.Post("/upload", UserHandler.UpdateProfilePhoto)
}

func (h UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	user := new(dto.UpdateUserRequest)
	if err := ctx.BodyParser(user); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(user); err != nil {
		return res.ValidationError(ctx, user.ToResponse(), err)
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	if err := h.UserUsecase.UpdateUser(userId, *user); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.UpdateSuccess, fiber.Map{
		"updated_data": user.ToResponse(),
	})
}

func (h UserHandler) UpdateProfilePhoto(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("photo")
	if err != nil {
		return res.BadRequest(ctx)
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	_err := h.UserUsecase.UpdateProfilePhoto(userId, file)
	if _err != nil {
		return res.Error(ctx, _err)
	}

	return res.SuccessResponse(ctx, res.UpdateProfilePhotoSuccess, nil)
}
