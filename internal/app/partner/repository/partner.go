package repository

import (
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type PartnerMySQLItf interface {
	Get(partner *entity.Partner) error
	Create(partner *entity.Partner) error
}

type PartnerMySQL struct {
	db *gorm.DB
}

func NewPartnerMySQL(db *gorm.DB) PartnerMySQLItf {
	return &PartnerMySQL{db}
}

func (r *PartnerMySQL) Get(partner *entity.Partner) error {
	return r.db.Find(partner).Error
}

func (r *PartnerMySQL) Create(partner *entity.Partner) error {
	return r.db.Create(partner).Error
}
