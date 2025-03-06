package entity

import (
	"github.com/google/uuid"
	"time"
)

type PaymentStatus string

const (
	Capture    PaymentStatus = "capture"
	Challenge  PaymentStatus = "challenge"
	Accept     PaymentStatus = "accept"
	Settlement PaymentStatus = "settlement"
	Deny       PaymentStatus = "deny"
	Cancel     PaymentStatus = "cancel"
	Expire     PaymentStatus = "expire"
	Pending    PaymentStatus = "pending"
)

type Status string

const (
	Finish    Status = "finish"
	Process   Status = "process"
	Cancelled Status = "cancelled"
)

type Transaction struct {
	ID                 uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID             uuid.UUID `gorm:"type:char(36);not null"`
	PaymentMethodID    uuid.UUID `gorm:"type:char(36);not null"`
	TransactionDetails []TransactionDetail
	PaymentStatus      PaymentStatus `gorm:"type:ENUM('capture','challenge','accept','settlement','deny','cancel','expire','pending');default:null"`
	Invoice            string        `gorm:"type:varchar(255);not null;uniqueIndex"`
	Total              float32       `gorm:"type:float(24);not null"`
	Status             Status        `gorm:"type:ENUM('finish','process','cancelled');default:null"`
	Note               string        `gorm:"type:text"`
	CreatedAt          time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt          time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
