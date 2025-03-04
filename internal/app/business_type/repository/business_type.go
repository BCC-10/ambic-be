package repository

import (
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type BusinessTypeMySQLItf interface {
	Get(businessType *entity.BusinessType) error
	Create(businessType *entity.BusinessType) error
}

type BusinessTypeMySQL struct {
	db *gorm.DB
}

func NewBusinessTypeMySQL(db *gorm.DB) BusinessTypeMySQLItf {
	return &BusinessTypeMySQL{db}
}

func (r BusinessTypeMySQL) Get(businessType *entity.BusinessType) error {
	return r.db.Debug().Preload("Partner").Find(&businessType).Error
}

func (r BusinessTypeMySQL) Create(businessType *entity.BusinessType) error {
	return r.db.Debug().Create(businessType).Error
}
