package dto

type GetTransactionDetailResponse struct {
	ProductID string `json:"product_id"`
	Qty       uint   `json:"qty"`
}

type CreateTransactionDetailRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Qty       int    `json:"qty" validate:"required,numeric,min=1"`
}
