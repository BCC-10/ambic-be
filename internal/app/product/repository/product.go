package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type ProductMySQLItf interface {
	Create(product *entity.Product) error
	Update(product *entity.Product) error
	Show(product *entity.Product, param dto.ProductParam) error
	Delete(product *entity.Product) error
}

type ProductMySQL struct {
	db *gorm.DB
}

func NewProductMySQL(db *gorm.DB) ProductMySQLItf {
	return &ProductMySQL{db}
}

func (r *ProductMySQL) Show(product *entity.Product, param dto.ProductParam) error {
	return r.db.Debug().Preload("Partner").First(&product, param).Error
}

func (r *ProductMySQL) Create(product *entity.Product) error {
	return r.db.Debug().Create(product).Error
}

func (r *ProductMySQL) Update(product *entity.Product) error {
	return r.db.Debug().Updates(product).Error
}

func (r *ProductMySQL) Delete(product *entity.Product) error {
	return r.db.Debug().Delete(product).Error
}
