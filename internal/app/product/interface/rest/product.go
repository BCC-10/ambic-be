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
	helper         helper.HelperIf
}

func NewProductHandler(routerGroup fiber.Router, productUsecase usecase.ProductUsecaseItf, validator *validator.Validate, m middleware.MiddlewareIf, helper helper.HelperIf) {
	ProductHandler := ProductHandler{
		ProductUsecase: productUsecase,
		validator:      validator,
		helper:         helper,
	}

	routerGroup = routerGroup.Group("/products")
	routerGroup.Post("/", m.Authentication, m.EnsurePartner, m.EnsureVerifiedPartner, ProductHandler.CreateProduct)
	routerGroup.Delete("/:id", m.Authentication, m.EnsurePartner, m.EnsureVerifiedPartner, ProductHandler.DeleteProduct)
	routerGroup.Patch("/:id/update", m.Authentication, m.EnsurePartner, m.EnsureVerifiedPartner, ProductHandler.UpdateProduct)
}

func (h ProductHandler) CreateProduct(ctx *fiber.Ctx) error {
	req := new(dto.CreateProductRequest)
	if err := h.helper.FormParser(ctx, req); err != nil {
		return res.BadRequest(ctx)
	}

	partnerId := ctx.Locals("partnerId").(uuid.UUID)

	if err := h.validator.Struct(req); err != nil {
		return res.ValidationError(ctx, nil, err)
	}

	if err := h.ProductUsecase.CreateProduct(partnerId, *req); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.CreateProductSuccess, nil)
}

func (h ProductHandler) UpdateProduct(ctx *fiber.Ctx) error {
	req := new(dto.UpdateProductRequest)
	if err := h.helper.FormParser(ctx, req); err != nil {
		return res.BadRequest(ctx)
	}

	productId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx, res.InvalidUUID)
	}

	if err := h.validator.Struct(req); err != nil {
		return res.ValidationError(ctx, nil, err)
	}

	partnerId := ctx.Locals("partnerId").(uuid.UUID)
	if err := h.ProductUsecase.UpdateProduct(productId, partnerId, *req); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.UpdateProductSuccess, nil)
}

func (h ProductHandler) DeleteProduct(ctx *fiber.Ctx) error {
	productId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx, res.InvalidUUID)
	}

	partnerId := ctx.Locals("partnerId").(uuid.UUID)
	if err := h.ProductUsecase.DeleteProduct(productId, partnerId); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.DeleteProductSuccess, nil)
}
