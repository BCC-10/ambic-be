package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type NotificationMySQLItf interface {
	GetByUserId(notif *[]entity.Notification, param dto.NotificationParam) error
	Create(notif *entity.Notification) error
}

type NotificationMySQL struct {
	db *gorm.DB
}

func NewNotificationMySQL(db *gorm.DB) NotificationMySQLItf {
	return &NotificationMySQL{db}
}

func (r *NotificationMySQL) GetByUserId(notif *[]entity.Notification, param dto.NotificationParam) error {
	return r.db.Debug().Find(notif, param).Error
}

func (r *NotificationMySQL) Create(notif *entity.Notification) error {
	return r.db.Debug().Create(notif).Error
}
