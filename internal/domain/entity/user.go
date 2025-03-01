package entity

import (
	"github.com/google/uuid"
	"time"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type User struct {
	ID         uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Username   string    `gorm:"type:varchar(255);unique;not null"`
	Email      string    `gorm:"type:varchar(255);unique;not null"`
	Phone      string    `gorm:"type:varchar(15);uniqueIndex;default:null"`
	Address    string    `gorm:"type:text;default:null"`
	BornDate   time.Time `gorm:"type:date;default:null"`
	Gender     Gender    `gorm:"type:ENUM('male','female');default:null"`
	Password   string    `gorm:"type:varchar(255)"`
	IsVerified bool      `gorm:"type:boolean;default:false"`
	PhotoURL   string    `gorm:"type:varchar(255);default:null"`
	CreatedAt  time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
