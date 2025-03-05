package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessType struct {
	ID       uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	Name     string    `gorm:"type:varchar(255);not null"`
	Partners []Partner
}

func (b *BusinessType) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewUUID()
	b.ID = id
	return
}
