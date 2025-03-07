package repository

import (
	"gorm.io/gorm"
)

type TransactionDetailMySQLItf interface {
}

type TransactionDetailMySQL struct {
	db *gorm.DB
}

func NewTransactionMySQL(db *gorm.DB) TransactionDetailMySQLItf {
	return &TransactionDetailMySQL{db}
}
