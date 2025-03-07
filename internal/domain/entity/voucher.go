package entity

import "github.com/google/uuid"

type Type string

const (
	Percentage Type = "percentage"
	Amount     Type = "amount"
)

type Voucher struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey"`
	Code     string    `gorm:"type:varchar(255);unique;not null"`
	Type     Type      `gorm:"type:ENUM('percentage','amount');default:null"`
	Discount int       `gorm:"type:int;not null"`
}
