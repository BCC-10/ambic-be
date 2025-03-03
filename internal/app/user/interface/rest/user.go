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

func NewUserHandler(routerGroup fiber.Router, userUsecase usecase.UserUsecaseItf, validator *validator.Validate, m middleware.MiddlewareIf) {
	UserHandler := UserHandler{
		Validator:   validator,
		UserUsecase: userUsecase,
	}

	routerGroup = routerGroup.Group("/users")
	routerGroup.Patch("/update", m.Authentication, UserHandler.UpdateUser)
}

func (h UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	req := new(dto.UpdateUserRequest)
	if ctx.Get("Content-Type") == "application/json" {
		if err := ctx.BodyParser(req); err != nil {
			return res.BadRequest(ctx)
		}
	} else {
		req = &dto.UpdateUserRequest{
			Name:     ctx.FormValue("name"),
			Phone:    ctx.FormValue("phone"),
			Address:  ctx.FormValue("address"),
			BornDate: ctx.FormValue("born_date"),
			Gender:   ctx.FormValue("gender"),
		}
	}

	if err := h.Validator.Struct(req); err != nil {
		return res.ErrValidationError(ctx, err)
	}

	file, err := ctx.FormFile("photo")
	if err == nil {
		req.Photo = file
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	if err := h.UserUsecase.UpdateUser(userId, *req); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.UpdateSuccess, fiber.Map{
		"user": req.ToResponse(),
	})
}
