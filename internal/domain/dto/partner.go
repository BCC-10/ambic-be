package dto

type RegisterPartnerRequest struct {
	Name      string  `form:"name" json:"name" validate:"required"`
	Type      string  `form:"type" json:"type" validate:"required"`
	Address   string  `form:"address" json:"address" validate:"required"`
	City      string  `form:"city" json:"city" validate:"required"`
	Longitude float64 `form:"longitude" json:"longitude" validate:"required,longitude"`
	Latitude  float64 `form:"latitude" json:"latitude" validate:"required,latitude"`
}

type RegisterPartnerResponse struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Address   string  `json:"address"`
	City      string  `json:"city"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (r *RegisterPartnerRequest) AsResponse() RegisterPartnerResponse {
	return RegisterPartnerResponse{
		Name:      r.Name,
		Type:      r.Type,
		Address:   r.Address,
		City:      r.City,
		Longitude: r.Longitude,
		Latitude:  r.Latitude,
	}
}
