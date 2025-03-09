package entity

import (
	"ambic/internal/domain/dto"
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
	ID            uuid.UUID      `gorm:"type:char(36);primaryKey"`
	Partner       Partner        `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	Ratings       []Rating       `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	Transactions  []Transaction  `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	Notifications []Notification `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	Name          string         `gorm:"type:varchar(255);default:null"`
	Username      string         `gorm:"type:varchar(255);unique;not null"`
	Email         string         `gorm:"type:varchar(255);unique;not null"`
	Phone         string         `gorm:"type:varchar(15);uniqueIndex;default:null"`
	Address       string         `gorm:"type:text;default:null"`
	BornDate      time.Time      `gorm:"type:date;default:null"`
	Gender        *Gender        `gorm:"type:ENUM('male','female');default:null"`
	Password      string         `gorm:"type:varchar(255)"`
	IsVerified    bool           `gorm:"type:boolean;default:false"`
	PhotoURL      string         `gorm:"type:varchar(255);not null"`
	CreatedAt     time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;autoUpdateTime"`
}

func (t *User) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	t.ID = id
	return
}

func (g Gender) String() string {
	return string(g)
}

func (t *User) ParseDTOGet() dto.GetUserResponse {
	getUserResponse := dto.GetUserResponse{
		ID:       t.ID.String(),
		Name:     t.Name,
		Username: t.Username,
		Email:    t.Email,
		Phone:    t.Phone,
		Address:  t.Address,
		Photo:    t.PhotoURL,
	}

	if t.Partner.ID != uuid.Nil {
		getUserResponse.Partner = t.Partner.ParseDTOGet()
	}

	if t.BornDate != (time.Time{}) {
		getUserResponse.BornDate = t.BornDate.String()
	}

	if t.Gender != nil {
		getUserResponse.Gender = t.Gender.String()
	}

	return getUserResponse
}
