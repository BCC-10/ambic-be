package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethod struct {
	ID           uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	Name         string    `gorm:"type:varchar(255)"`
	Transactions []Transaction
}

func (u *PaymentMethod) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewUUID()
	u.ID = id
	return
}
