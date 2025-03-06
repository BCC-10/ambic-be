package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type BusinessTypeMySQLItf interface {
	Show(businessType *entity.BusinessType, param dto.BusinessTypeParam) error
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

func (r BusinessTypeMySQL) Show(businessType *entity.BusinessType, param dto.BusinessTypeParam) error {
	return r.db.Debug().First(businessType, param).Error
}
