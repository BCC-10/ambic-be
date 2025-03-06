package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type TransactionMySQLItf interface {
	Get(transaction *[]entity.Transaction, param dto.TransactionParam) error
	Create(transaction *entity.Transaction) error
	Update(transaction *entity.Transaction) error
}

type TransactionMySQL struct {
	db *gorm.DB
}

func NewTransactionMySQL(db *gorm.DB) TransactionMySQLItf {
	return &TransactionMySQL{db}
}

func (r *TransactionMySQL) Get(transaction *[]entity.Transaction, param dto.TransactionParam) error {
	return r.db.Debug().Find(transaction, param).Error
}

func (r *TransactionMySQL) Update(transaction *entity.Transaction) error {
	return r.db.Debug().Updates(transaction).Error
}

func (r *TransactionMySQL) Create(transaction *entity.Transaction) error {
	return r.db.Debug().Create(transaction).Error
}
