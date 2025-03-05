package entity

import (
	"ambic/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Partner struct {
	ID             uuid.UUID `gorm:"type:varchar(36);primary_key"`
	UserID         uuid.UUID `gorm:"type:varchar(36);not null;uniqueIndex"`
	BusinessTypeID uuid.UUID `gorm:"type:varchar(36);not null"`
	Products       []Product
	Name           string    `gorm:"type:varchar(255);not null"`
	Type           string    `gorm:"type:varchar(255);not null"`
	Address        string    `gorm:"type:varchar(255);not null"`
	City           string    `gorm:"type:varchar(255);not null"`
	Longitude      float64   `gorm:"type:float;not null"`
	Latitude       float64   `gorm:"type:float;not null"`
	Instagram      string    `gorm:"type:varchar(255);not null"`
	IsVerified     bool      `gorm:"type:boolean;default:false"`
	PhotoURL       string    `gorm:"type:varchar(255);default:null"`
	CreatedAt      time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

func (p *Partner) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	p.ID = id
	return
}

func (p *Partner) ParseDTOGet() dto.GetPartnerResponse {
	return dto.GetPartnerResponse{
		ID:        p.ID.String(),
		Name:      p.Name,
		Type:      p.Type,
		Address:   p.Address,
		City:      p.City,
		Instagram: p.Instagram,
		Longitude: p.Longitude,
		Latitude:  p.Latitude,
		Photo:     p.PhotoURL,
	}
}
