package rest

import (
	"ambic/internal/app/partner/usecase"
	"ambic/internal/domain/dto"
	"ambic/internal/infra/helper"
	res "ambic/internal/infra/response"
	"ambic/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PartnerHandler struct {
	Validator      *validator.Validate
	PartnerUsecase usecase.PartnerUsecaseItf
	helper         helper.HelperIf
}

func NewPartnerHandler(routerGroup fiber.Router, partnerUsecase usecase.PartnerUsecaseItf, validator *validator.Validate, m middleware.MiddlewareIf, helper helper.HelperIf) {
	PartnerHandler := PartnerHandler{
		PartnerUsecase: partnerUsecase,
		Validator:      validator,
		helper:         helper,
	}

	routerGroup = routerGroup.Group("/partners")
	routerGroup.Get("/:id", m.Authentication, m.EnsurePartner, PartnerHandler.ShowPartner)
	routerGroup.Get("/:id/products", m.Authentication, m.EnsurePartner, m.EnsureVerifiedPartner, PartnerHandler.GetProducts)
	routerGroup.Post("/", m.Authentication, m.EnsureNotPartner, PartnerHandler.RegisterPartner)
	routerGroup.Post("/verification", m.Authentication, PartnerHandler.VerifyPartner)
	routerGroup.Patch("/", m.Authentication, m.EnsurePartner, m.EnsureVerifiedPartner, PartnerHandler.UpdatePhoto)
}

func (h *PartnerHandler) RegisterPartner(ctx *fiber.Ctx) error {
	data := new(dto.RegisterPartnerRequest)
	if err := h.helper.FormParser(ctx, data); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return res.ValidationError(ctx, nil, err)
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	token, err := h.PartnerUsecase.RegisterPartner(userId, *data)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.PartnerRegisterSuccess, fiber.Map{
		"new_token": token,
	})
}

func (h *PartnerHandler) VerifyPartner(ctx *fiber.Ctx) error {
	data := new(dto.VerifyPartnerRequest)
	if err := ctx.BodyParser(&data); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return res.ValidationError(ctx, nil, err)
	}

	if err := h.PartnerUsecase.VerifyPartner(*data); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.PartnerVerifySuccess, nil)
}

func (h *PartnerHandler) GetProducts(ctx *fiber.Ctx) error {
	query := new(dto.GetPartnerProductsQuery)
	if err := ctx.QueryParser(query); err != nil {
		return res.BadRequest(ctx)
	}

	partnerId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx, res.InvalidUUID)
	}

	products, _err := h.PartnerUsecase.GetProducts(partnerId, *query)
	if _err != nil {
		return res.Error(ctx, _err)
	}

	return res.SuccessResponse(ctx, res.GetProductSuccess, fiber.Map{
		"products": products,
	})
}

func (h *PartnerHandler) ShowPartner(ctx *fiber.Ctx) error {
	partnerId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx, res.InvalidUUID)
	}

	partner, _err := h.PartnerUsecase.ShowPartner(partnerId)
	if _err != nil {
		return res.Error(ctx, _err)
	}

	return res.SuccessResponse(ctx, res.GetPartnerSuccess, fiber.Map{
		"partner": partner,
	})
}

func (h *PartnerHandler) UpdatePhoto(ctx *fiber.Ctx) error {
	data := new(dto.UpdatePhotoRequest)
	if err := h.helper.FormParser(ctx, data); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return res.ValidationError(ctx, nil, err)
	}

	partnerId := ctx.Locals("partnerId").(uuid.UUID)

	if err := h.PartnerUsecase.UpdatePhoto(partnerId, *data); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.UpdatePartnerPhotoSuccess, nil)
}
