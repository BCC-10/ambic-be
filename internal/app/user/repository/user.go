package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type UserMySQLItf interface {
	Login(user *entity.User, param dto.LoginRequest) error
	Show(user *entity.User, param dto.UserParam) error
	Create(user *entity.User) error
	Verify(user *entity.User) error
	Update(user *entity.User) error
}

type UserMySQL struct {
	db *gorm.DB
}

func NewUserMySQL(db *gorm.DB) UserMySQLItf {
	return &UserMySQL{db}
}

func (r UserMySQL) Login(user *entity.User, param dto.LoginRequest) error {
	return r.db.Debug().Where("email = ? OR username = ?", param.Identifier, param.Identifier).Preload("Partner").First(&user).Error
}

func (r UserMySQL) Show(user *entity.User, param dto.UserParam) error {
	return r.db.Debug().Preload("Partner").First(&user, param).Error
}

func (r UserMySQL) Create(user *entity.User) error {
	return r.db.Debug().Create(user).Error
}

func (r UserMySQL) Update(user *entity.User) error {
	return r.db.Debug().Updates(user).Error
}

func (r UserMySQL) Verify(user *entity.User) error {
	return r.db.Debug().Model(&user).Where("email = ?", user.Email).Update("is_verified", true).Error
}
