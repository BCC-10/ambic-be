package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type PartnerMySQLItf interface {
	Get(partner *[]entity.Partner, param dto.PartnerParam) error
	Show(partner *entity.Partner, param dto.PartnerParam) error
	ShowWithTransactions(partner *entity.Partner, param dto.PartnerParam) error
	Create(tx *gorm.DB, partner *entity.Partner) error
	Update(tx *gorm.DB, partner *entity.Partner) error
}

type PartnerMySQL struct {
	db *gorm.DB
}

func NewPartnerMySQL(db *gorm.DB) PartnerMySQLItf {
	return &PartnerMySQL{db}
}

func (r *PartnerMySQL) Get(partner *[]entity.Partner, param dto.PartnerParam) error {
	return r.db.Debug().Find(partner, param).Error
}

func (r *PartnerMySQL) Create(tx *gorm.DB, partner *entity.Partner) error {
	return tx.Debug().Create(partner).Error
}

func (r *PartnerMySQL) Update(tx *gorm.DB, partner *entity.Partner) error {
	return tx.Debug().Updates(partner).Error
}

func (r *PartnerMySQL) Show(partner *entity.Partner, param dto.PartnerParam) error {
	return r.db.Debug().Preload("BusinessType").First(partner, param).Error
}

func (r *PartnerMySQL) ShowWithTransactions(partner *entity.Partner, param dto.PartnerParam) error {
	return r.db.Debug().Preload("Transactions").First(partner, param).Error
}
