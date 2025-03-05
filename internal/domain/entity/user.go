package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type User struct {
	ID         uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	Partner    Partner
	Name       string    `gorm:"type:varchar(255)"`
	Username   string    `gorm:"type:varchar(255);unique;not null"`
	Email      string    `gorm:"type:varchar(255);unique;not null"`
	Phone      string    `gorm:"type:varchar(15);uniqueIndex;default:null"`
	Address    string    `gorm:"type:text;default:null"`
	BornDate   time.Time `gorm:"type:date;default:null"`
	Gender     *Gender   `gorm:"type:ENUM('male','female');default:null"`
	Password   string    `gorm:"type:varchar(255)"`
	IsVerified bool      `gorm:"type:boolean;default:false"`
	PhotoURL   string    `gorm:"type:varchar(255);not null"`
	CreatedAt  time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
