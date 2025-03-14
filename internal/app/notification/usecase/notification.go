package usecase

import (
	"ambic/internal/app/notification/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/helper"
	res "ambic/internal/infra/response"
	"github.com/google/uuid"
)

type NotificationUsecaseItf interface {
	GetByUserId(userId uuid.UUID, pagination dto.PaginationRequest) ([]dto.GetNotificationResponse, *dto.PaginationResponse, *res.Err)
}

type NotificationUsecase struct {
	env                    *env.Env
	helper                 helper.HelperIf
	NotificationRepository repository.NotificationMySQLItf
}

func NewNotificationUsecase(env *env.Env, notificationRepository repository.NotificationMySQLItf, helper helper.HelperIf) NotificationUsecaseItf {
	return &NotificationUsecase{
		env:                    env,
		helper:                 helper,
		NotificationRepository: notificationRepository,
	}
}

func (u *NotificationUsecase) GetByUserId(userId uuid.UUID, pagination dto.PaginationRequest) ([]dto.GetNotificationResponse, *dto.PaginationResponse, *res.Err) {
	pagination = u.helper.CreatePagination(pagination)

	notifications := new([]entity.Notification)
	totalNotifications, err := u.NotificationRepository.GetByUserId(notifications, dto.NotificationParam{UserID: userId}, pagination)
	if err != nil {
		return nil, nil, res.ErrInternalServer()
	}

	var response []dto.GetNotificationResponse
	for _, notification := range *notifications {
		response = append(response, notification.ParseDTOGet())
	}

	pg := u.helper.CalculatePagination(pagination, totalNotifications)

	return response, &pg, nil
}
