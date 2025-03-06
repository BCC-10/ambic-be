package rest

import (
	RatingUsecase "ambic/internal/app/rating/usecase"
	"ambic/internal/domain/dto"
	"ambic/internal/infra/helper"
	res "ambic/internal/infra/response"
	"ambic/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RatingHandler struct {
	RatingUsecase RatingUsecase.RatingUsecaseItf
	validator     *validator.Validate
	helper        helper.HelperIf
}

func NewRatingHandler(routerGroup fiber.Router, ratingUsecase RatingUsecase.RatingUsecaseItf, validator *validator.Validate, m middleware.MiddlewareIf, helper helper.HelperIf) {
	RatingHandler := RatingHandler{
		RatingUsecase: ratingUsecase,
		validator:     validator,
		helper:        helper,
	}

	routerGroup = routerGroup.Group("/ratings")
	routerGroup.Get("/", m.Authentication, RatingHandler.Get)
	routerGroup.Get("/:id", m.Authentication, RatingHandler.Show)
	routerGroup.Post("/", m.Authentication, RatingHandler.Create)
	routerGroup.Patch("/:id/update", m.Authentication, RatingHandler.Update)
	routerGroup.Delete("/:id", m.Authentication, RatingHandler.Delete)
}

func (h *RatingHandler) Get(ctx *fiber.Ctx) error {
	ratings, err := h.RatingUsecase.Get()
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.GetRatingSuccess, ratings)
}

func (h *RatingHandler) Show(ctx *fiber.Ctx) error {
	ratingId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx)
	}

	rating, _err := h.RatingUsecase.Show(ratingId)
	if _err != nil {
		return res.Error(ctx, _err)
	}

	return res.SuccessResponse(ctx, res.GetRatingSuccess, rating)
}

func (h *RatingHandler) Create(ctx *fiber.Ctx) error {
	req := new(dto.CreateRatingRequest)
	if err := h.helper.FormParser(ctx, req); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.validator.Struct(req); err != nil {
		return res.ValidationError(ctx, nil, err)
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	if err := h.RatingUsecase.Create(userId, *req); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.RatingCreateSuccess, nil)
}

func (h *RatingHandler) Update(ctx *fiber.Ctx) error {
	req := new(dto.UpdateRatingRequest)
	if err := h.helper.FormParser(ctx, req); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.validator.Struct(req); err != nil {
		return res.ValidationError(ctx, nil, err)
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	ratingId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.RatingUsecase.Update(userId, ratingId, *req); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.RatingUpdateSuccess, nil)
}

func (h *RatingHandler) Delete(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(uuid.UUID)

	ratingId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.RatingUsecase.Delete(userId, ratingId); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.RatingDeleteSuccess, nil)
}
