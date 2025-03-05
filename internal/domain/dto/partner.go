package dto

import "github.com/google/uuid"

type RegisterPartnerRequest struct {
	Name      string  `form:"name" validate:"required"`
	Type      string  `form:"type" validate:"required"`
	Address   string  `form:"address" validate:"required"`
	City      string  `form:"city" validate:"required"`
	Instagram string  `form:"instagram" validate:"required"`
	Longitude float64 `form:"longitude" validate:"required,longitude"`
	Latitude  float64 `form:"latitude" validate:"required,latitude"`
}

type VerifyPartnerRequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

type GetPartnerProductsQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type PartnerParam struct {
	ID    uuid.UUID
	Email string
}
