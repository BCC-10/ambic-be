package mysql

import (
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(entity.User{}, entity.BusinessType{}, entity.Partner{}, entity.Rating{}, entity.Product{}, entity.Payment{}, entity.Transaction{}, entity.TransactionDetail{}, entity.Notification{})
}
