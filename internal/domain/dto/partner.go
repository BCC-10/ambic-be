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
	Longitude      float64               `form:"longitude" validate:"required,longitude"`
	Latitude       float64               `form:"latitude" validate:"required,latitude"`
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

type GetPartnerProductsQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type GetPartnerResponse struct {
	ID             string  `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	BusinessTypeID string  `json:"type,omitempty"`
	Address        string  `json:"address,omitempty"`
	City           string  `json:"city,omitempty"`
	Instagram      string  `json:"instagram,omitempty"`
	Longitude      float64 `json:"longitude,omitempty"`
	Latitude       float64 `json:"latitude,omitempty"`
	Photo          string  `json:"photo,omitempty"`
}

type PartnerParam struct {
	ID    uuid.UUID
	Email string
}
