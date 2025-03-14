package dto

import (
	"github.com/google/uuid"
	"time"
)

type GetTransactionByUserIdAndByStatusRequest struct {
	Status string `query:"status" validate:"omitempty,oneof='waiting for payment' 'process' 'finish' 'cancelled by system'"`
	Limit  int    `query:"limit"`
	Page   int    `query:"page"`
}

type GetTransactionResponse struct {
	ID         string               `json:"id"`
	User       GetUserResponse      `json:"user"`
	Payment    GetPaymentResponse   `json:"payment,omitempty"`
	Invoice    string               `json:"invoice"`
	Total      float32              `json:"total"`
	Status     string               `json:"status"`
	Note       string               `json:"note"`
	Datetime   time.Time            `json:"datetime"`
	Items      []GetProductResponse `json:"items"`
	PaymentURL string               `json:"payment_url"`
}

type CreateTransactionRequest struct {
	Note               string                           `json:"note"`
	PartnerID          string                           `json:"partner_id" validate:"required,uuid"`
	TransactionDetails []CreateTransactionDetailRequest `json:"items" validate:"required"`
}

type UpdateTransactionStatusRequest struct {
	Status string `json:"status" validate:"required,oneof='finish' 'cancelled by system'"`
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
