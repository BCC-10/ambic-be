package repository

import (
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type BusinessTypeMySQLItf interface {
	Create(businessType *entity.BusinessType) error
}

type BusinessTypeMySQL struct {
	db *gorm.DB
}

func NewBusinessTypeMySQL(db *gorm.DB) BusinessTypeMySQLItf {
	return &BusinessTypeMySQL{db}
}

func (r BusinessTypeMySQL) Create(businessType *entity.BusinessType) error {
	return r.db.Debug().Create(businessType).Error
}
