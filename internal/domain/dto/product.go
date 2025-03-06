package dto

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type GetProductResponse struct {
	ID           string  `json:"id"`
	PartnerID    string  `json:"partner_id,omitempty"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	InitialPrice float32 `json:"initial_price"`
	FinalPrice   float32 `json:"final_price"`
	Stock        int     `json:"stock"`
	PickupTime   string  `json:"pickup_time"`
	PhotoURL     string  `json:"photo"`
}

type CreateProductRequest struct {
	Name         string                `form:"name" validate:"required"`
	Description  string                `form:"description" validate:"required"`
	InitialPrice float32               `form:"initial_price" validate:"required,numeric"`
	FinalPrice   float32               `form:"final_price" validate:"required,numeric"`
	Stock        int                   `form:"stock" validate:"required,numeric"`
	PickupTime   string                `form:"pickup_time" validate:"required,datetime=2006-01-02 15:04:05"`
	Photo        *multipart.FileHeader `json:"photo" form:"photo" validate:"required"`
}

type UpdateProductRequest struct {
	Name         string                `form:"name"`
	Description  string                `form:"description"`
	InitialPrice float32               `form:"initial_price" validate:"omitempty,numeric"`
	FinalPrice   float32               `form:"final_price" validate:"omitempty,numeric"`
	Stock        int                   `form:"stock" validate:"omitempty,numeric"`
	PickupTime   string                `form:"pickup_time" validate:"omitempty,datetime=2006-01-02 15:04:05"`
	Photo        *multipart.FileHeader `form:"photo"`
}

type ProductParam struct {
	Id        uuid.UUID
	PartnerId uuid.UUID
	Name      string
}
