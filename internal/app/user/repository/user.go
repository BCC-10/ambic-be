package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type UserMySQLItf interface {
	Get(user *entity.User, param dto.UserParam) error
	Create(user *entity.User) error
	Activate(user *entity.User) error
}

type UserMySQL struct {
	db *gorm.DB
}

func NewUserMySQL(db *gorm.DB) UserMySQLItf {
	return &UserMySQL{db}
}

func (r UserMySQL) Get(user *entity.User, param dto.UserParam) error {
	return r.db.Debug().First(&user, param).Error
}

func (r UserMySQL) Create(user *entity.User) error {
	return r.db.Debug().Create(user).Error
}

func (r UserMySQL) Activate(user *entity.User) error {
	return r.db.Debug().Model(&user).Where("email = ?", user.Email).Update("is_active", true).Error
}
