package repository

import (
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type PaymentMethodMySQLItf interface {
	Create(paymentMethod *entity.PaymentMethod) error
}

type PaymentMethodMySQL struct {
	db *gorm.DB
}

func NewPaymentMethodMySQL(db *gorm.DB) PaymentMethodMySQLItf {
	return &PaymentMethodMySQL{db}
}

func (r PaymentMethodMySQL) Create(paymentMethod *entity.PaymentMethod) error {
	return r.db.Debug().Create(paymentMethod).Error
}
