package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	ID                uuid.UUID `gorm:"type:char(36);primaryKey"`
	TransactionID     uuid.UUID `gorm:"type:char(36);not null"`
	ReferenceID       string    `gorm:"type:varchar(255);not null"`
	TransactionStatus string    `gorm:"type:varchar(255);not null"`
	StatusMessage     string    `gorm:"type:varchar(255);not null"`
	PaymentType       string    `gorm:"type:varchar(255);not null"`
	FraudStatus       string    `gorm:"type:varchar(255);not null"`
	TransactionTime   time.Time `gorm:"type:timestamp"`
	SettlementTime    time.Time `gorm:"type:timestamp"`
}

func (u *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewUUID()
	u.ID = id
	return
}
