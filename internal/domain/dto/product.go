package dto

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type FilterProductRequest struct {
	Name   string  `query:"name"`
	Lat    float64 `query:"lat" validate:"required"`
	Long   float64 `query:"long" validate:"required"`
	Radius float64 `query:"radius" validate:"required"`
	Limit  int     `query:"limit"`
	Page   int     `query:"page"`
	Offset int
}

type GetProductResponse struct {
	ID            string  `json:"id"`
	PartnerID     string  `json:"partner_id,omitempty"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	InitialPrice  float32 `json:"initial_price"`
	FinalPrice    float32 `json:"final_price"`
	Stock         uint    `json:"stock"`
	PickupTime    string  `json:"pickup_time"`
	EndPickupTime string  `json:"end_pickup_time"`
	PhotoURL      string  `json:"photo"`
	Star          float32 `json:"star,omitempty"`
	CountRating   int     `json:"count_rating,omitempty"`
	Distance      float64 `json:"distance,omitempty"`
}

type CreateProductRequest struct {
	Name          string                `form:"name" validate:"required"`
	Description   string                `form:"description" validate:"required"`
	InitialPrice  float32               `form:"initial_price" validate:"required,numeric,min=1"`
	FinalPrice    float32               `form:"final_price" validate:"required,numeric,min=1"`
	Stock         int                   `form:"stock" validate:"required,numeric"`
	PickupTime    string                `form:"pickup_time" validate:"required,datetime=2006-01-02 15:04:05"`
	EndPickupTime string                `form:"end_pickup_time" validate:"required,datetime=2006-01-02 15:04:05"`
	Photo         *multipart.FileHeader `form:"photo" validate:"required"`
}

type UpdateProductRequest struct {
	Name          string                `form:"name"`
	Description   string                `form:"description"`
	InitialPrice  float32               `form:"initial_price" validate:"omitempty,numeric,min=1"`
	FinalPrice    float32               `form:"final_price" validate:"omitempty,numeric,min=1"`
	Stock         int                   `form:"stock" validate:"omitempty,numeric"`
	PickupTime    string                `form:"pickup_time" validate:"omitempty,datetime=2006-01-02 15:04:05"`
	EndPickupTime string                `form:"end_pickup_time" validate:"omitempty,datetime=2006-01-02 15:04:05"`
	Photo         *multipart.FileHeader `form:"photo"`
}

type ProductParam struct {
	ID        uuid.UUID
	PartnerId uuid.UUID
	Name      string
}
