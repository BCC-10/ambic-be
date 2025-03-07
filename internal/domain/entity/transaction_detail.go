package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionDetail struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey"`
	TransactionID uuid.UUID `gorm:"type:char(36);not null"`
	ProductID     uuid.UUID `gorm:"type:char(36);not null"`
	Qty           uint      `gorm:"type:int;not null"`
}

func (d *TransactionDetail) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewUUID()
	d.ID = id
	return
}
