package usecase

import (
	"ambic/internal/app/notification/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
	"github.com/google/uuid"
)

type NotificationUsecaseItf interface {
	GetByUserId(userId uuid.UUID, pagination dto.PaginationRequest) ([]dto.GetNotificationResponse, *res.Err)
}

type NotificationUsecase struct {
	env                    *env.Env
	NotificationRepository repository.NotificationMySQLItf
}

func NewNotificationUsecase(env *env.Env, notificationRepository repository.NotificationMySQLItf) NotificationUsecaseItf {
	return &NotificationUsecase{
		env:                    env,
		NotificationRepository: notificationRepository,
	}
}

func (u *NotificationUsecase) GetByUserId(userId uuid.UUID, pagination dto.PaginationRequest) ([]dto.GetNotificationResponse, *res.Err) {
	if pagination.Limit < 1 {
		pagination.Limit = u.env.DefaultPaginationLimit
	}

	if pagination.Page < 1 {
		pagination.Page = u.env.DefaultPaginationPage
	}

	pagination.Offset = (pagination.Page - 1) * pagination.Limit

	notifications := new([]entity.Notification)
	if err := u.NotificationRepository.GetByUserId(notifications, dto.NotificationParam{UserID: userId}, pagination); err != nil {
		return nil, res.ErrInternalServer()
	}

	var response []dto.GetNotificationResponse
	for _, notification := range *notifications {
		response = append(response, notification.ParseDTOGet())
	}

	return response, nil
}
