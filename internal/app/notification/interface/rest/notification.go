package rest

import (
	"ambic/internal/app/notification/usecase"
	"ambic/internal/domain/dto"
	res "ambic/internal/infra/response"
	"ambic/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	NotificationUsecase usecase.NotificationUsecaseItf
}

func NewNotificationHandler(routerGroup fiber.Router, notificationUsecase usecase.NotificationUsecaseItf, m middleware.MiddlewareIf) {
	NotificationHandler := NotificationHandler{
		NotificationUsecase: notificationUsecase,
	}

	routerGroup = routerGroup.Group("/notifications")
	routerGroup.Get("/", m.Authentication, NotificationHandler.GetByUserId)
}

func (h *NotificationHandler) GetByUserId(ctx *fiber.Ctx) error {
	pagination := new(dto.PaginationRequest)
	if err := ctx.QueryParser(pagination); err != nil {
		return res.BadRequest(ctx)
	}

	userId := ctx.Locals("userId").(uuid.UUID)
	notifications, err := h.NotificationUsecase.GetByUserId(userId, *pagination)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.GetNotificationSuccess, notifications)
}
