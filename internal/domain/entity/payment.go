package entity

import (
	"ambic/internal/domain/dto"
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
	CreatedAt         time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt         time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

func (u *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewUUID()
	u.ID = id
	return
}

func (u *Payment) ParseDTOGet() dto.GetPaymentResponse {
	return dto.GetPaymentResponse{
		ID:                u.ID,
		TransactionID:     u.TransactionID,
		ReferenceID:       u.ReferenceID,
		TransactionStatus: u.TransactionStatus,
		StatusMessage:     u.StatusMessage,
		PaymentType:       u.PaymentType,
		FraudStatus:       u.FraudStatus,
		TransactionTime:   u.TransactionTime.String(),
		SettlementTime:    u.SettlementTime.String(),
	}
}
