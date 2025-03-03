package dto

type RegisterPartnerRequest struct {
	Name      string  `form:"name" json:"name" validate:"required"`
	Type      string  `form:"type" json:"type" validate:"required"`
	Address   string  `form:"address" json:"address" validate:"required"`
	City      string  `form:"city" json:"city" validate:"required"`
	Longitude float64 `form:"longitude" json:"longitude" validate:"required,longitude"`
	Latitude  float64 `form:"latitude" json:"latitude" validate:"required,latitude"`
}
