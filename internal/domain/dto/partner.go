package dto

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type RegisterPartnerRequest struct {
	Name           string                `form:"name" validate:"required"`
	Address        string                `form:"address" validate:"required"`
	City           string                `form:"city" validate:"required"`
	Instagram      string                `form:"instagram" validate:"required"`
	PlaceID        string                `form:"place_id" validate:"required"`
	BusinessTypeID string                `form:"business_type_id" validate:"required,uuid"`
	Photo          *multipart.FileHeader `form:"photo"`
}

type VerifyPartnerRequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

type UpdatePhotoRequest struct {
	Photo *multipart.FileHeader `form:"photo" validate:"required"`
}

type GetPartnerStatisticResponse struct {
	TotalRatings      int64   `json:"total_ratings"`
	TotalProducts     int64   `json:"total_products"`
	TotalTransactions int64   `json:"total_transactions"`
	TotalRevenue      float32 `json:"total_revenue"`
}

type GetPartnerResponse struct {
	PlaceID      string  `json:"place_id,omitempty"`
	ID           string  `json:"id,omitempty"`
	Name         string  `json:"name,omitempty"`
	BusinessType string  `json:"business_type,omitempty"`
	Address      string  `json:"address,omitempty"`
	City         string  `json:"city,omitempty"`
	Instagram    string  `json:"instagram,omitempty"`
	Longitude    float64 `json:"longitude,omitempty"`
	Latitude     float64 `json:"latitude,omitempty"`
	Photo        string  `json:"photo,omitempty"`
}

type GetPartnerTransactionRequest struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Status string `query:"status" validate:"omitempty,oneof='waiting for payment' 'process' 'finish' 'cancelled by system'"`
	Offset int
}

type PartnerParam struct {
	ID         uuid.UUID
	IsVerified bool
	Email      string
}
