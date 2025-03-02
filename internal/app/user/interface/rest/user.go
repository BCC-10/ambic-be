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
}

func (h UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	data := new(dto.UpdateUserRequest)
	if ctx.Get("Content-Type") == "application/json" {
		if err := ctx.BodyParser(data); err != nil {
			return res.BadRequest(ctx)
		}
	} else {
		data = &dto.UpdateUserRequest{
			Name:     ctx.FormValue("name"),
			Phone:    ctx.FormValue("phone"),
			Address:  ctx.FormValue("address"),
			BornDate: ctx.FormValue("born_date"),
			Gender:   ctx.FormValue("gender"),
		}
	}

	err := h.Validator.Struct(data)
	if err != nil {
		return res.ErrValidationError(ctx, err)
	}

	file, err := ctx.FormFile("photo")
	if err == nil {
		data.Photo = file
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	_err := h.UserUsecase.UpdateUser(userId, *data)
	if _err != nil {
		return res.Error(ctx, _err)
	}

	return res.SuccessResponse(ctx, res.UpdateSuccess, fiber.Map{
		"data": data.ToResponse(),
	})
}
