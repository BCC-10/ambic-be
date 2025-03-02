package entity

type SocialMedia struct {
	ID       uint   `gorm:"primary_key"`
	Media    string `gorm:"type:varchar(255);not null"`
	Username string `gorm:"type:varchar(255);not null"`
}
