package entity

import (
	"ambic/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	ID                uuid.UUID `gorm:"type:char(36);primaryKey"`
	TransactionID     uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	OrderID           string    `gorm:"type:varchar(255)"`
	ReferenceID       string    `gorm:"type:varchar(255)"`
	MerchantID        string    `gorm:"type:varchar(255)"`
	Issuer            string    `gorm:"type:varchar(255)"`
	Currency          string    `gorm:"type:varchar(255)"`
	GrossAmount       float32   `gorm:"type:float(24)"`
	Acquirer          string    `gorm:"type:varchar(255)"`
	TransactionStatus string    `gorm:"type:varchar(255)"`
	StatusMessage     string    `gorm:"type:varchar(255)"`
	PaymentType       string    `gorm:"type:varchar(255)"`
	FraudStatus       string    `gorm:"type:varchar(255)"`
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
		ID:                u.ID.String(),
		TransactionID:     u.TransactionID.String(),
		OrderID:           u.OrderID,
		ReferenceID:       u.ReferenceID,
		TransactionStatus: u.TransactionStatus,
		StatusMessage:     u.StatusMessage,
		PaymentType:       u.PaymentType,
		FraudStatus:       u.FraudStatus,
		TransactionTime:   u.TransactionTime.String(),
		SettlementTime:    u.SettlementTime.String(),
	}
}
