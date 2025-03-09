package entity

import (
	"ambic/internal/domain/dto"
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	WaitingForPayment  Status = "waiting for payment"
	Finish             Status = "finish"
	Process            Status = "process"
	CancelledBySystem  Status = "cancelled by system"
	CancelledByUser    Status = "cancelled by user"
	CancelledByPartner Status = "cancelled by partner"
)

type Transaction struct {
	ID                 uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID             uuid.UUID `gorm:"type:char(36);not null"`
	PartnerID          uuid.UUID `gorm:"type:char(36)"`
	Payment            Payment
	TransactionDetails []TransactionDetail
	Invoice            string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Total              float32   `gorm:"type:float(24);not null"`
	Status             Status    `gorm:"type:ENUM('waiting for payment','finish','process','cancelled by user','cancelled by partner','cancelled by system', 'waiting for confirmation');default:null"`
	Note               string    `gorm:"type:text"`
	CreatedAt          time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

func (t *Transaction) ParseDTOGet() dto.GetTransactionResponse {
	products := make([]dto.GetProductResponse, len(t.TransactionDetails))
	for i, detail := range t.TransactionDetails {
		products[i] = detail.Product.ParseDTOGet()
	}

	res := dto.GetTransactionResponse{
		ID:       t.ID.String(),
		UserID:   t.UserID.String(),
		Invoice:  t.Invoice,
		Total:    t.Total,
		Status:   string(t.Status),
		Note:     t.Note,
		Datetime: t.UpdatedAt,
		Items:    products,
	}

	if t.Payment.ID != uuid.Nil {
		res.Payment = t.Payment.ParseDTOGet()
	}

	return res
}
