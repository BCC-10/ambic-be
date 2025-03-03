package dto

type RegisterPartnerRequest struct {
	Name      string  `form:"name" validate:"required"`
	Type      string  `form:"type" validate:"required"`
	Address   string  `form:"address" validate:"required"`
	City      string  `form:"city" validate:"required"`
	Longitude float64 `form:"longitude" validate:"required,longitude"`
	Latitude  float64 `form:"latitude" validate:"required,latitude"`
}
