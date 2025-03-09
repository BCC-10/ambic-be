package dto

type Pagination struct {
	Page   int `query:"page"`
	Limit  int `query:"limit"`
	Offset int
}
