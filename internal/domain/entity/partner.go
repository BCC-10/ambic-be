package entity

import (
	"ambic/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Partner struct {
	ID             uuid.UUID     `gorm:"type:char(36);primary_key"`
	UserID         uuid.UUID     `gorm:"type:char(36);uniqueIndex"`
	BusinessTypeID uuid.UUID     `gorm:"type:char(36);not null"`
	Products       []Product     `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	Transactions   []Transaction `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	BusinessType   BusinessType
	Name           string    `gorm:"type:varchar(255);not null"`
	Address        string    `gorm:"type:varchar(255);not null"`
	City           string    `gorm:"type:varchar(255);not null"`
	Longitude      float64   `gorm:"type:float;not null"`
	Latitude       float64   `gorm:"type:float;not null"`
	PlaceID        string    `gorm:"type:varchar(255);not null"`
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
		ID:           p.ID.String(),
		PlaceID:      p.PlaceID,
		Name:         p.Name,
		BusinessType: p.BusinessType.Name,
		Address:      p.Address,
		City:         p.City,
		Instagram:    p.Instagram,
		Longitude:    p.Longitude,
		Latitude:     p.Latitude,
		Photo:        p.PhotoURL,
	}
}
