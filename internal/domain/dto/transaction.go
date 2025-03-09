package dto

import (
	"github.com/google/uuid"
	"time"
)

type GetTransactionByUserIdAndByStatusRequest struct {
	Status string `query:"status" validate:"omitempty,oneof='waiting for payment' 'waiting for confirmation' 'process' 'finish' 'cancelled by system' 'cancelled by user' 'cancelled by partner'"`
	Limit  int    `query:"limit"`
	Page   int    `query:"page"`
}

type GetTransactionResponse struct {
	ID       string               `json:"id"`
	UserID   string               `json:"user_id"`
	Payment  GetPaymentResponse   `json:"payment,omitempty"`
	Invoice  string               `json:"invoice"`
	Total    float32              `json:"total"`
	Status   string               `json:"status"`
	Note     string               `json:"note"`
	Datetime time.Time            `json:"datetime"`
	Items    []GetProductResponse `json:"items"`
}

type CreateTransactionRequest struct {
	Note               string                           `json:"note"`
	PartnerID          string                           `json:"partner_id" validate:"required,uuid"`
	TransactionDetails []CreateTransactionDetailRequest `json:"items" validate:"required"`
}

type UpdateTransactionStatusRequest struct {
	Status string `json:"status" validate:"required,oneof='finish' 'cancelled by system' 'cancelled by user' 'cancelled by partner'"`
}

type RequestSnap struct {
	TransactionID      string
	OrderID            string
	Amount             int64
	User               GetUserResponse
	TransactionDetails []TransactionDetail
}

type TransactionParam struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PartnerID uuid.UUID
	ProductID uuid.UUID
	Status    string
}
