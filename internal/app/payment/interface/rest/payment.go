package rest

import (
	PaymentUsecase "ambic/internal/app/payment/usecase"
	"ambic/internal/domain/dto"
	res "ambic/internal/infra/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	Validator      *validator.Validate
	PaymentUsecase PaymentUsecase.PaymentUsecaseItf
}

func NewPaymentHandler(routerGroup fiber.Router, paymentUsecase PaymentUsecase.PaymentUsecaseItf, validator *validator.Validate) {
	PaymentHandler := PaymentHandler{
		Validator:      validator,
		PaymentUsecase: paymentUsecase,
	}

	routerGroup = routerGroup.Group("/payments")
	routerGroup.Post("/notification", PaymentHandler.Notification)
}

func (h PaymentHandler) Notification(ctx *fiber.Ctx) error {
	req := new(dto.NotificationPayment)
	if err := ctx.BodyParser(req); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(req); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.PaymentUsecase.ProcessPayment(req); err != nil {
		return res.Error(ctx, err)
	}

	return nil
}
