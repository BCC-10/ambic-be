package entity

import (
	"ambic/internal/domain/dto"
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	WaitingForPayment Status = "waiting for payment"
	Finish            Status = "finish"
	Process           Status = "process"
	Cancelled         Status = "cancelled"
)

type Transaction struct {
	ID                 uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID             uuid.UUID `gorm:"type:char(36);not null"`
	Payment            Payment
	TransactionDetails []TransactionDetail
	Invoice            string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Total              float32   `gorm:"type:float(24);not null"`
	Status             Status    `gorm:"type:ENUM('waiting for payment','finish','process','cancelled');default:null"`
	Note               string    `gorm:"type:text"`
	CreatedAt          time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

func (t *Transaction) ParseDTOGet() dto.GetTransactionResponse {
	res := dto.GetTransactionResponse{
		ID:        t.ID.String(),
		Invoice:   t.Invoice,
		Total:     t.Total,
		Status:    string(t.Status),
		Note:      t.Note,
		UpdatedAt: t.UpdatedAt,
	}

	if t.Payment.ID != uuid.Nil {
		res.Payment = t.Payment.ParseDTOGet()
	}

	return res
}
