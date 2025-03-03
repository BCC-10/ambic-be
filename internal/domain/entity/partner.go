package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Partner struct {
	ID         uuid.UUID `gorm:"type:varchar(36);primary_key"`
	Products   []Product
	UserID     uuid.UUID `gorm:"type:varchar(36);not null;"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Type       string    `gorm:"type:varchar(255);not null"`
	Address    string    `gorm:"type:varchar(255);not null"`
	City       string    `gorm:"type:varchar(255);not null"`
	Longitude  float64   `gorm:"type:float;not null"`
	Latitude   float64   `gorm:"type:float;not null"`
	IsVerified bool      `gorm:"type:boolean;default:false"`
	CreatedAt  time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

func (p *Partner) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
