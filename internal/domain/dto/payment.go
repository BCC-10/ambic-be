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
