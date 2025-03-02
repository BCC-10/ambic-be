package entity

type Product struct {
	ID           uint `gorm:"primaryKey"`
	PartnerID    uint `gorm:"type:int;not null"`
	Partner      Partner
	Name         string  `gorm:"type:varchar(255);not null"`
	Description  string  `gorm:"type:text;not null;"`
	InitialPrice float32 `gorm:"type:float;not null"`
	FinalPrice   float32 `gorm:"type:float;not null"`
	Stock        int     `gorm:"type:int;not null"`
	PickupTime   string  `gorm:"type:varchar(30);not null"`
	PhotoURL     string  `gorm:"type:varchar(255)"`
	CreatedAt    string  `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt    string  `gorm:"type:timestamp;autoUpdateTime"`
}
