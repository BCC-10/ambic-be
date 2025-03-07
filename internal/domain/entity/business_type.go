package entity

import (
	"ambic/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessType struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name     string    `gorm:"type:varchar(255);not null"`
	Partners []Partner
}

func (b *BusinessType) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewUUID()
	b.ID = id
	return
}

func (b *BusinessType) ParseDTOGet() dto.GetBusinessTypeResponse {
	return dto.GetBusinessTypeResponse{
		ID:   b.ID,
		Name: b.Name,
	}
}
