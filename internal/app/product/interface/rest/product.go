package rest

import (
	"ambic/internal/app/product/usecase"
	"ambic/internal/domain/dto"
	"ambic/internal/infra/helper"
	res "ambic/internal/infra/response"
	"ambic/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	ProductUsecase usecase.ProductUsecaseItf
	validator      *validator.Validate
}

func NewProductHandler(routerGroup fiber.Router, productUsecase usecase.ProductUsecaseItf, validator *validator.Validate, m middleware.MiddlewareIf) {
	ProductHandler := ProductHandler{
		ProductUsecase: productUsecase,
		validator:      validator,
	}

	routerGroup = routerGroup.Group("/products")
	routerGroup.Post("/create", m.Authentication, m.EnsurePartner, m.EnsureVerifiedPartner, ProductHandler.CreateProduct)
	routerGroup.Patch("/:id/update", m.Authentication, m.EnsurePartner, m.EnsureVerifiedPartner, ProductHandler.UpdateProduct)
}

func (h ProductHandler) CreateProduct(ctx *fiber.Ctx) error {
	req := new(dto.CreateProductRequest)
	if err := helper.ParseForm(ctx, req); err != nil {
		return res.BadRequest(ctx)
	}

	userId := ctx.Locals("userId").(uuid.UUID)

	if err := h.validator.Struct(req); err != nil {
		return res.ValidationError(ctx, req.ToResponse(), err)
	}

	if err := h.ProductUsecase.CreateProduct(userId, *req); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.ProductCreateSuccess, req.ToResponse())
}

func (h ProductHandler) UpdateProduct(ctx *fiber.Ctx) error {
	req := new(dto.UpdateProductRequest)
	if err := helper.ParseForm(ctx, req); err != nil {
		return res.BadRequest(ctx)
	}

	productId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.validator.Struct(req); err != nil {
		return res.ValidationError(ctx, nil, err)
	}

	partnerId := ctx.Locals("partnerId").(uuid.UUID)
	if err := h.ProductUsecase.UpdateProduct(productId, partnerId, *req); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.ProductUpdateSuccess, nil)
}
