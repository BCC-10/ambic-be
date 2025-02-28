package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Username   string    `gorm:"type:varchar(255);unique;not null"`
	Email      string    `gorm:"type:varchar(255);unique;not null"`
	Phone      string    `gorm:"type:varchar(15);uniqueIndex;default:null"`
	Password   string    `gorm:"type:varchar(255)"`
	IsVerified bool      `gorm:"type:boolean;default:false"`
	CreatedAt  time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
