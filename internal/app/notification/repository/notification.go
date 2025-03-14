package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type NotificationMySQLItf interface {
	GetByUserId(notif *[]entity.Notification, param dto.NotificationParam, pagination dto.PaginationRequest) (int64, error)
	Create(tx *gorm.DB, notif *entity.Notification) error
}

type NotificationMySQL struct {
	db *gorm.DB
}

func NewNotificationMySQL(db *gorm.DB) NotificationMySQLItf {
	return &NotificationMySQL{db}
}

func (r *NotificationMySQL) GetByUserId(notif *[]entity.Notification, param dto.NotificationParam, pagination dto.PaginationRequest) (int64, error) {
	result := r.db.Debug().Limit(pagination.Limit).Offset(pagination.Offset).Order("created_at desc").Find(notif, param)

	var count int64

	result.Count(&count)

	return count, result.Error
}

func (r *NotificationMySQL) Create(tx *gorm.DB, notif *entity.Notification) error {
	return tx.Debug().Create(notif).Error
}
