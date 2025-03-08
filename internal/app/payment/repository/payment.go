package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type PaymentMySQLItf interface {
	Create(payment *entity.Payment) error
	Show(payment *entity.Payment, param dto.PaymentParam) error
}

type PaymentMySQL struct {
	db *gorm.DB
}

func NewPaymentMySQL(db *gorm.DB) PaymentMySQLItf {
	return &PaymentMySQL{db}
}

func (r *PaymentMySQL) Create(payment *entity.Payment) error {
	return r.db.Debug().Create(payment).Error
}

func (r *PaymentMySQL) Show(payment *entity.Payment, param dto.PaymentParam) error {
	return r.db.Debug().First(payment, param).Error
}
