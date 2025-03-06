package rest

import (
	"ambic/internal/app/transaction/usecase"
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

	routerGroup = routerGroup.Group("/transactions")
	routerGroup.Get("/", m.Authentication, TransactionHandler.GetByUser)
}

func (h *TransactionHandler) GetByUser(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(uuid.UUID)
	transactions, err := h.TransactionUsecase.GetByUserID(userId)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "suc,", fiber.Map{
		"transactions": transactions,
	})
}
