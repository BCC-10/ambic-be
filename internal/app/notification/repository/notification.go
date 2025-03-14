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

func (r *NotificationMySQL) GetByUserId(notification *[]entity.Notification, param dto.NotificationParam, pagination dto.PaginationRequest) (int64, error) {
	query := r.db.Debug()

	var count int64
	if err := query.Model(&notification).Count(&count).Error; err != nil {
		return 0, err
	}

	query.Limit(pagination.Limit).Offset(pagination.Offset).Order("created_at desc").Find(notification, param)

	return count, query.Error
}

func (r *NotificationMySQL) Create(tx *gorm.DB, notif *entity.Notification) error {
	return tx.Debug().Create(notif).Error
}
