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
	"log"
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
	routerGroup.Get("/location", m.Authentication, PartnerHandler.AutocompleteLocation)
	routerGroup.Get("/:id/products", m.Authentication, m.EnsurePartner, m.EnsureVerifiedPartner, PartnerHandler.GetProducts)
	routerGroup.Get("/:id/transactions", PartnerHandler.GetTransactions)
	routerGroup.Get("/:id/statistics", m.Authentication, m.EnsurePartner, m.EnsureVerifiedPartner, PartnerHandler.GetStatistics)
	routerGroup.Get("/:id", m.Authentication, m.EnsurePartner, PartnerHandler.ShowPartner)
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

	token, err := h.PartnerUsecase.VerifyPartner(*data)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.PartnerVerificationSuccess, fiber.Map{
		"new_token": token,
	})
}

func (h *PartnerHandler) GetProducts(ctx *fiber.Ctx) error {
	pagination := new(dto.Pagination)
	if err := ctx.QueryParser(pagination); err != nil {
		return res.BadRequest(ctx)
	}

	log.Println(pagination)

	partnerId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx, res.InvalidUUID)
	}

	products, _err := h.PartnerUsecase.GetProducts(partnerId, *pagination)
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

func (h *PartnerHandler) AutocompleteLocation(ctx *fiber.Ctx) error {
	req := new(dto.LocationRequest)
	if err := ctx.QueryParser(req); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(req); err != nil {
		return res.ValidationError(ctx, nil, err)
	}

	data, err := h.PartnerUsecase.AutocompleteLocation(*req)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.GetAutoCompleteSuccess, fiber.Map{
		"locations": data,
	})
}

func (h *PartnerHandler) GetStatistics(ctx *fiber.Ctx) error {
	partnerId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx, res.InvalidUUID)
	}

	statistic, _err := h.PartnerUsecase.GetStatistics(partnerId)
	if _err != nil {
		return res.Error(ctx, _err)
	}

	return res.SuccessResponse(ctx, res.GetPartnerStatisticsSuccess, fiber.Map{
		"statistic": statistic,
	})
}

func (h *PartnerHandler) GetTransactions(ctx *fiber.Ctx) error {
	pagination := new(dto.Pagination)
	if err := ctx.QueryParser(pagination); err != nil {
		return res.BadRequest(ctx)
	}

	partnerId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx, res.InvalidUUID)
	}

	transactions, _err := h.PartnerUsecase.GetTransactions(partnerId, *pagination)
	if _err != nil {
		return res.Error(ctx, _err)
	}

	return res.SuccessResponse(ctx, res.GetTransactionSuccess, fiber.Map{
		"transactions": transactions,
	})
}
