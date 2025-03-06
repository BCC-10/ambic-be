package dto

import (
	"github.com/google/uuid"
	"time"
)

type GetTransactionResponse struct {
	ID        uuid.UUID          `json:"id"`
	Payment   GetPaymentResponse `json:"payment"`
	Invoice   string             `json:"invoice"`
	Total     float32            `json:"total"`
	Status    string             `json:"status"`
	Note      string             `json:"note"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type TransactionParam struct {
	ID     uuid.UUID
	UserID uuid.UUID
}
