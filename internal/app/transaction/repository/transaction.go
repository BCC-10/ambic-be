package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionMySQLItf interface {
	Get(transaction *[]entity.Transaction, param dto.TransactionParam, pagination dto.PaginationRequest) (int64, error)
	Show(transaction *entity.Transaction, param dto.TransactionParam) error
	Create(tx *gorm.DB, transaction *entity.Transaction) error
	Update(tx *gorm.DB, transaction *entity.Transaction) error
	CheckIsUserHasPurchasedProduct(param dto.TransactionParam) bool
}

type TransactionMySQL struct {
	db *gorm.DB
}

func NewTransactionMySQL(db *gorm.DB) TransactionMySQLItf {
	return &TransactionMySQL{db}
}

func (r *TransactionMySQL) Get(transaction *[]entity.Transaction, param dto.TransactionParam, pagination dto.PaginationRequest) (int64, error) {
	query := r.db.Debug().Preload(clause.Associations).Preload("TransactionDetails.Product")

	var count int64
	if err := query.Model(&transaction).Count(&count).Error; err != nil {
		return 0, err
	}

	query.Limit(pagination.Limit).Offset(pagination.Offset).Order("created_at desc").Find(transaction, param)

	return count, query.Error
}

func (r *TransactionMySQL) CheckIsUserHasPurchasedProduct(param dto.TransactionParam) bool {
	var count int64
	r.db.Table("transactions").
		Select("COUNT(*)").
		Joins("JOIN transaction_details ON transactions.id = transaction_details.transaction_id").
		Where("transactions.user_id = ? AND transaction_details.product_id = ? AND transactions.status = ?", param.UserID, param.ProductID, entity.Finish).
		Count(&count)

	return count > 0
}

func (r *TransactionMySQL) Update(tx *gorm.DB, transaction *entity.Transaction) error {
	return tx.Debug().Preload("TransactionDetails").Updates(transaction).Error
}

func (r *TransactionMySQL) Create(tx *gorm.DB, transaction *entity.Transaction) error {
	return tx.Debug().Create(transaction).Error
}

func (r *TransactionMySQL) Show(transaction *entity.Transaction, param dto.TransactionParam) error {
	return r.db.Debug().Preload(clause.Associations).Preload("TransactionDetails.Product").First(transaction, param).Error
}
