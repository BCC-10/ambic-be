package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewV7()
	return
}
