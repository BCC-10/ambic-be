package dto

type GetProductsResponse struct {
	ID           uint    `json:"id"`
	PartnerID    uint    `json:"partner_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	InitialPrice float32 `json:"initial_price"`
	FinalPrice   float32 `json:"final_price"`
	Stock        int     `json:"stock"`
	PickupTime   string  `json:"pickup_time"`
	PhotoURL     string  `json:"photo"`
}
