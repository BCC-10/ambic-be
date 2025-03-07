package dto

type NotificationPayment struct {
	OrderID           string `json:"order_id" validate:"required"`
	ReferenceID       string `json:"reference_id" validate:"required"`
	TransactionStatus string `json:"transaction_status" validate:"required"`
	StatusMessage     string `json:"status_message" validate:"required"`
	PaymentType       string `json:"payment_type" validate:"required"`
	FraudStatus       string `json:"fraud_status" validate:"required"`
	TransactionTime   string `json:"transaction_time" validate:"required"`
	SettlementTime    string `json:"settlement_time" validate:"required"`
}

type GetPaymentResponse struct {
	ID                string `json:"id,omitempty"`
	TransactionID     string `json:"transaction_id,omitempty"`
	OrderID           string `json:"order_id,omitempty"`
	ReferenceID       string `json:"reference_id,omitempty"`
	TransactionStatus string `json:"transaction_status,omitempty"`
	StatusMessage     string `json:"status_message,omitempty"`
	PaymentType       string `json:"payment_type,omitempty"`
	FraudStatus       string `json:"fraud_status,omitempty"`
	TransactionTime   string `json:"transaction_time,omitempty"`
	SettlementTime    string `json:"settlement_time,omitempty"`
}
