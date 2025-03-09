package dto

import (
	"github.com/google/uuid"
	"time"
)

type GetTransactionResponse struct {
	ID      string             `json:"id"`
	Payment GetPaymentResponse `json:"payment,omitempty"`
	Invoice string             `json:"invoice"`
	Total   float32            `json:"total"`
	Status  string             `json:"status"`
	Note    string             `json:"note"`
	Date    time.Time          `json:"date"`
}

type ShowTransactionResponse struct {
	ID      string               `json:"id"`
	Payment GetPaymentResponse   `json:"payment,omitempty"`
	Invoice string               `json:"invoice"`
	Total   float32              `json:"total"`
	Status  string               `json:"status"`
	Note    string               `json:"note"`
	Date    time.Time            `json:"date"`
	Items   []GetProductResponse `json:"items"`
}

type CreateTransactionRequest struct {
	Note               string                           `json:"note"`
	PartnerID          string                           `json:"partner_id" validate:"required,uuid"`
	TransactionDetails []CreateTransactionDetailRequest `json:"items" validate:"required"`
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
	ProductID uuid.UUID
}
