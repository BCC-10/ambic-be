package dto

type LocationRequest struct {
	Query  string  `json:"query" validate:"required"`
	Lat    float64 `json:"lat"`
	Long   float64 `json:"long"`
	Radius float64 `json:"radius"`
}

type LocationResponse struct {
	Name    string `json:"name"`
	PlaceID string `json:"place_id"`
}

type PlaceDetails struct {
	Name    string  `json:"name"`
	PlaceId string  `json:"place_id"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
}
