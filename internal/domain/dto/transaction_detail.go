package dto

type TransactionDetail struct {
	MerchantName string
	ProductID    string
	Product      GetProductResponse
	Qty          uint
}

type GetTransactionDetailResponse struct {
	ProductID string             `json:"product_id"`
	Product   GetProductResponse `json:"product"`
	Qty       uint               `json:"qty"`
}

type CreateTransactionDetailRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Qty       int    `json:"qty" validate:"required,numeric,min=1"`
}
