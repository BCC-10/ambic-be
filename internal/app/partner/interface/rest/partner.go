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

	routerGroup = routerGroup.Group("/partners", m.Authentication)
	routerGroup.Get("/:id/products", m.EnsurePartner, m.EnsureVerifiedPartner, PartnerHandler.GetProducts)
	routerGroup.Get("/:id/transactions", m.EnsurePartner, m.EnsureVerifiedPartner, PartnerHandler.GetTransactions)
	routerGroup.Get("/:id/statistics", m.EnsurePartner, PartnerHandler.GetStatistics)
	routerGroup.Get("/:id", m.EnsurePartner, PartnerHandler.ShowPartner)
	routerGroup.Post("/", m.EnsureNotPartner, PartnerHandler.RegisterPartner)
	routerGroup.Post("/verification", PartnerHandler.VerifyPartner)
	routerGroup.Patch("/", m.EnsurePartner, m.EnsureVerifiedPartner, PartnerHandler.UpdatePhoto)
}

func (h *PartnerHandler) RegisterPartner(ctx *fiber.Ctx) error {
	data := new(dto.RegisterPartnerRequest)
	if err := h.helper.FormParser(ctx, data); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return res.ValidationError(ctx, err)
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
		return res.ValidationError(ctx, err)
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
	pagination := new(dto.PaginationRequest)
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
		return res.ValidationError(ctx, err)
	}

	partnerId := ctx.Locals("partnerId").(uuid.UUID)

	if err := h.PartnerUsecase.UpdatePhoto(partnerId, *data); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.UpdatePartnerPhotoSuccess, nil)
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
	pagination := new(dto.PaginationRequest)
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
