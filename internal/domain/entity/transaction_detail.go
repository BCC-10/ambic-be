package entity

import (
	"ambic/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionDetail struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey"`
	TransactionID uuid.UUID `gorm:"type:char(36);not null"`
	ProductID     uuid.UUID `gorm:"type:char(36)"`
	Product       Product
	Qty           uint `gorm:"type:int;not null"`
}

func (d *TransactionDetail) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewUUID()
	d.ID = id
	return
}

func (d *TransactionDetail) ParseDTOGet() dto.GetTransactionDetailResponse {
	return dto.GetTransactionDetailResponse{
		ProductID: d.ProductID.String(),
		Product:   d.Product.ParseDTOGet(nil),
		Qty:       d.Qty,
	}
}
