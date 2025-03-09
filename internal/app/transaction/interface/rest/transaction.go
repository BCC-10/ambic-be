package rest

import (
	"ambic/internal/app/transaction/usecase"
	"ambic/internal/domain/dto"
	res "ambic/internal/infra/response"
	"ambic/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	TransactionUsecase usecase.TransactionUsecaseItf
	Validator          *validator.Validate
}

func NewTransactionHandler(routerGroup fiber.Router, transactionUsecase usecase.TransactionUsecaseItf, validator *validator.Validate, m middleware.MiddlewareIf) {
	TransactionHandler := TransactionHandler{
		TransactionUsecase: transactionUsecase,
		Validator:          validator,
	}

	routerGroup = routerGroup.Group("/transactions", m.Authentication)
	routerGroup.Get("/", TransactionHandler.GetByLoggedInUser)
	routerGroup.Get("/:id", TransactionHandler.Show)
	routerGroup.Post("/", m.EnsurePartner, m.EnsureVerifiedPartner, TransactionHandler.Create)
}

func (h *TransactionHandler) GetByLoggedInUser(ctx *fiber.Ctx) error {
	pagination := new(dto.PaginationRequest)
	if err := ctx.QueryParser(pagination); err != nil {
		return res.BadRequest(ctx)
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	transactions, err := h.TransactionUsecase.GetByUserID(userId, *pagination)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.GetTransactionSuccess, fiber.Map{
		"transactions": transactions,
	})
}

func (h *TransactionHandler) Create(ctx *fiber.Ctx) error {
	req := new(dto.CreateTransactionRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(req); err != nil {
		return res.ValidationError(ctx, nil, err)
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	paymentURL, err := h.TransactionUsecase.Create(userId, req)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.CreateTransactionSuccess, fiber.Map{
		"payment_url": paymentURL,
	})
}

func (h *TransactionHandler) Show(ctx *fiber.Ctx) error {
	transactionId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.BadRequest(ctx, res.InvalidUUID)
	}

	transactionDetails, _err := h.TransactionUsecase.Show(transactionId)
	if _err != nil {
		return res.Error(ctx, _err)
	}

	return res.SuccessResponse(ctx, res.GetTransactionSuccess, transactionDetails)
}
