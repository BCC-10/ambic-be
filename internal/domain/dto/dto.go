package dto

type PaginationRequest struct {
	Page   int `query:"page"`
	Limit  int `query:"limit"`
	Offset int
}

type Location struct {
	Lat  float64
	Long float64
}
