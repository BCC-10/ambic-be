package repository

import (
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type UserMySQLItf interface {
	Create(user *entity.User) error
}

type UserMySQL struct {
	db *gorm.DB
}

func NewUserMySQL(db *gorm.DB) UserMySQLItf {
	return &UserMySQL{db}
}

func (r UserMySQL) Create(user *entity.User) error {
	return r.db.Debug().Create(user).Error
}
