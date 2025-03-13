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

type PartnerRegistrationTelegramMessage struct {
	UserID               string
	UserName             string
	UserUsername         string
	UserEmail            string
	UserPhone            string
	UserAddress          string
	UserGender           string
	UserRegisteredAt     string
	PartnerID            string
	BusinessType         string
	BusinessName         string
	BusinessAddress      string
	BusinessCity         string
	BusinessGmaps        string
	BusinessInstagram    string
	BusinessPhoto        string
	BusinessRegisteredAt string
}
