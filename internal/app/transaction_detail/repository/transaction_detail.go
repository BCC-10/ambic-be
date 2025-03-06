package repository

import (
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type TransactionDetailMySQLItf interface {
	Create(transactionDetail *entity.TransactionDetail) error
}

type TransactionDetailMySQL struct {
	db *gorm.DB
}

func NewTransactionMySQL(db *gorm.DB) TransactionDetailMySQLItf {
	return &TransactionDetailMySQL{db}
}

func (r *TransactionDetailMySQL) Create(transactionDetail *entity.TransactionDetail) error {
	return r.db.Debug().Create(transactionDetail).Error
}
