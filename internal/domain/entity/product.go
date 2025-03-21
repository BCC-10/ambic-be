package entity

import (
	"ambic/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID                 uuid.UUID `gorm:"type:char(36);primaryKey"`
	PartnerID          uuid.UUID `gorm:"type:char(36)"`
	Partner            Partner
	Ratings            []Rating            `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	TransactionDetails []TransactionDetail `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	Name               string              `gorm:"type:varchar(255);not null;uniqueIndex:idx_partner_product"`
	Description        string              `gorm:"type:text;not null;"`
	InitialPrice       float32             `gorm:"type:float;not null"`
	FinalPrice         float32             `gorm:"type:float;not null"`
	Stock              uint                `gorm:"type:int;not null"`
	PickupTime         time.Time           `gorm:"type:timestamp;not null"`
	EndPickupTime      time.Time           `gorm:"type:timestamp;not null"`
	PhotoURL           string              `gorm:"type:varchar(255)"`
	CreatedAt          time.Time           `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt          time.Time           `gorm:"type:timestamp;autoUpdateTime"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	p.ID = id
	return
}

func (p *Product) ParseDTOGet(distance *float64) dto.GetProductResponse {
	var star float32
	var countRating int
	if len(p.Ratings) > 0 {
		var totalRating float32
		for _, r := range p.Ratings {
			totalRating += float32(r.Star)
			countRating++
		}
		star = totalRating / float32(len(p.Ratings))
	}

	res := dto.GetProductResponse{
		ID:            p.ID.String(),
		PartnerID:     p.PartnerID.String(),
		Name:          p.Name,
		Description:   p.Description,
		InitialPrice:  p.InitialPrice,
		FinalPrice:    p.FinalPrice,
		Stock:         p.Stock,
		PickupTime:    p.PickupTime.String(),
		EndPickupTime: p.EndPickupTime.String(),
		PhotoURL:      p.PhotoURL,
		Star:          star,
	}

	if distance != nil {
		res.Distance = *distance
	}

	if countRating > 0 {
		res.CountRating = countRating
	}

	return res
}
