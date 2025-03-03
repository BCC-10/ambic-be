package entity

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	PartnerID    uuid.UUID `gorm:"type:varchar(36);not null"`
	Partner      Partner
	Name         string    `gorm:"type:varchar(255);not null"`
	Description  string    `gorm:"type:text;not null;"`
	InitialPrice float32   `gorm:"type:float;not null"`
	FinalPrice   float32   `gorm:"type:float;not null"`
	Stock        int       `gorm:"type:int;not null"`
	PickupTime   string    `gorm:"type:varchar(30);not null"`
	PhotoURL     string    `gorm:"type:varchar(255)"`
	CreatedAt    time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

func (p *Product) BeforeCreate() (err error) {
	p.ID = uuid.New()
	return
}
