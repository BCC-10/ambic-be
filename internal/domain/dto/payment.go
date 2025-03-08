package dto

import "github.com/google/uuid"

type NotificationPayment struct {
	OrderID           string `json:"order_id"`
	TransactionID     string `json:"custom_field1"`
	ReferenceID       string `json:"reference_id"`
	MerchantID        string `json:"merchant_id"`
	Issuer            string `json:"issuer"`
	GrossAmount       string `json:"gross_amount"`
	Currency          string `json:"currency"`
	Acquirer          string `json:"acquirer"`
	TransactionStatus string `json:"transaction_status"`
	StatusMessage     string `json:"status_message"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
	TransactionTime   string `json:"transaction_time"`
	SettlementTime    string `json:"settlement_time"`
}

type GetPaymentResponse struct {
	ID                string  `json:"id,omitempty"`
	OrderID           string  `json:"order_id,omitempty"`
	TransactionID     string  `json:"custom_field1,omitempty"`
	ReferenceID       string  `json:"reference_id,omitempty"`
	MerchantID        string  `json:"merchant_id,omitempty"`
	Issuer            string  `json:"issuer,omitempty"`
	GrossAmount       float32 `json:"gross_amount,omitempty"`
	Currency          string  `json:"currency,omitempty"`
	Acquirer          string  `json:"acquirer,omitempty"`
	TransactionStatus string  `json:"transaction_status,omitempty"`
	StatusMessage     string  `json:"status_message,omitempty"`
	PaymentType       string  `json:"payment_type,omitempty"`
	FraudStatus       string  `json:"fraud_status,omitempty"`
	TransactionTime   string  `json:"transaction_time,omitempty"`
	SettlementTime    string  `json:"settlement_time,omitempty"`
}

type PaymentParam struct {
	ID            uuid.UUID
	TransactionID uuid.UUID
}
