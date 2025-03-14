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

	routerGroup = routerGroup.Group("/ratings", m.Authentication)
	routerGroup.Get("/", RatingHandler.Get)
	routerGroup.Post("/", RatingHandler.Create)
	routerGroup.Get("/:id", RatingHandler.Show)
	routerGroup.Patch("/:id", RatingHandler.Update)
	routerGroup.Delete("/:id", RatingHandler.Delete)
}

func (h *RatingHandler) Get(ctx *fiber.Ctx) error {
	pagination := new(dto.PaginationRequest)
	if err := ctx.QueryParser(pagination); err != nil {
		return res.BadRequest(ctx)
	}

	req := new(dto.GetRatingRequest)
	if err := ctx.QueryParser(req); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.validator.Struct(req); err != nil {
		return res.ValidationError(ctx, err)
	}

	ratings, pg, err := h.RatingUsecase.Get(*req, *pagination)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.GetRatingSuccess, fiber.Map{
		"ratings": ratings,
		"pg":      pg,
	})
}

func (h *RatingHandler) Show(ctx *fiber.Ctx) error {
	req := new(dto.ShowRatingRequest)
	if err := ctx.ParamsParser(req); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.validator.Struct(req); err != nil {
		return res.ValidationError(ctx, err)
	}

	rating, _err := h.RatingUsecase.Show(*req)
	if _err != nil {
		return res.Error(ctx, _err)
	}

	return res.SuccessResponse(ctx, res.GetRatingSuccess, fiber.Map{
		"rating": rating,
	})
}

func (h *RatingHandler) Create(ctx *fiber.Ctx) error {
	req := new(dto.CreateRatingRequest)
	if err := h.helper.FormParser(ctx, req); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.validator.Struct(req); err != nil {
		return res.ValidationError(ctx, err)
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	if err := h.RatingUsecase.Create(userId, *req); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.CreateRatingSuccess, nil)
}

func (h *RatingHandler) Update(ctx *fiber.Ctx) error {
	req := new(dto.UpdateRatingRequest)
	if err := h.helper.FormParser(ctx, req); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.validator.Struct(req); err != nil {
		return res.ValidationError(ctx, err)
	}

	param := new(dto.UpdateRatingParam)
	if err := ctx.ParamsParser(param); err != nil {
		return res.BadRequest(ctx, res.InvalidUUID)
	}

	userId := ctx.Locals("userId").(uuid.UUID)

	if err := h.RatingUsecase.Update(userId, *param, *req); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.UpdateRatingSuccess, nil)
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
