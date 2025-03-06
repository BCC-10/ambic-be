package dto

type GetTransactionDetailResponse struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	ProductID     string `json:"product_id"`
	Qty           int    `json:"qty"`
}
