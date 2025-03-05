package dto

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type RegisterPartnerRequest struct {
	Name      string                `form:"name" validate:"required"`
	Type      string                `form:"type" validate:"required"`
	Address   string                `form:"address" validate:"required"`
	City      string                `form:"city" validate:"required"`
	Instagram string                `form:"instagram" validate:"required"`
	Longitude float64               `form:"longitude" validate:"required,longitude"`
	Latitude  float64               `form:"latitude" validate:"required,latitude"`
	Photo     *multipart.FileHeader `form:"photo"`
}

type VerifyPartnerRequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

type UpdatePhotoRequest struct {
	Photo *multipart.FileHeader `form:"photo" validate:"required"`
}

type GetPartnerProductsQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type GetPartnerResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Address   string  `json:"address"`
	City      string  `json:"city"`
	Instagram string  `json:"instagram"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Photo     string  `json:"photo"`
}

type PartnerParam struct {
	ID    uuid.UUID
	Email string
}
