package repository

import (
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type TransactionMySQLItf interface {
	Update(transaction *entity.Transaction) error
}

type TransactionMySQL struct {
	db *gorm.DB
}

func NewTransactionMySQL(db *gorm.DB) TransactionMySQLItf {
	return &TransactionMySQL{db}
}

func (r *TransactionMySQL) Update(transaction *entity.Transaction) error {
	return r.db.Debug().Updates(transaction).Error
}
