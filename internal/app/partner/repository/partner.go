package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type PartnerMySQLItf interface {
	GetProducts(partner *entity.Partner, param dto.PartnerParam, limit int, offset int) error
	Show(partner *entity.Partner, param dto.PartnerParam) error
	Create(partner *entity.Partner) error
	Update(partner *entity.Partner) error
}

type PartnerMySQL struct {
	db *gorm.DB
}

func NewPartnerMySQL(db *gorm.DB) PartnerMySQLItf {
	return &PartnerMySQL{db}
}

func (r *PartnerMySQL) GetProducts(partner *entity.Partner, param dto.PartnerParam, limit int, offset int) error {
	return r.db.Debug().Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit).Offset(offset)
	}).First(partner, param).Error
}

func (r *PartnerMySQL) Create(partner *entity.Partner) error {
	return r.db.Debug().Create(partner).Error
}

func (r *PartnerMySQL) Update(partner *entity.Partner) error {
	return r.db.Debug().Updates(partner).Error
}

func (r *PartnerMySQL) Show(partner *entity.Partner, param dto.PartnerParam) error {
	return r.db.Debug().First(partner, param).Error
}
