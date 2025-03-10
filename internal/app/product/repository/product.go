package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ProductMySQLItf interface {
	Create(product *entity.Product) error
	Show(product *entity.Product, param dto.ProductParam) error
	Delete(product *entity.Product) error
	GetByPartnerId(products *[]entity.Product, param dto.ProductParam, pagination dto.PaginationRequest) error
	Update(tx *gorm.DB, product *entity.Product) error
	GetTotalProductsByPartnerId(id uuid.UUID) (int64, error)
	Filter(products *[]entity.Product, param dto.ProductParam, pagination dto.PaginationRequest) error
}

type ProductMySQL struct {
	db *gorm.DB
}

func NewProductMySQL(db *gorm.DB) ProductMySQLItf {
	return &ProductMySQL{db}
}

func (r *ProductMySQL) Show(product *entity.Product, param dto.ProductParam) error {
	return r.db.Debug().Preload("Ratings").Preload("Partner").First(&product, param).Error
}

func (r *ProductMySQL) Create(product *entity.Product) error {
	return r.db.Debug().Create(product).Error
}

func (r *ProductMySQL) Delete(product *entity.Product) error {
	return r.db.Debug().Delete(product).Error
}

func (r *ProductMySQL) GetByPartnerId(product *[]entity.Product, param dto.ProductParam, pagination dto.PaginationRequest) error {
	return r.db.Debug().Preload("Ratings").Limit(pagination.Limit).Offset(pagination.Offset).Find(product, param).Error
}

func (r *ProductMySQL) Update(tx *gorm.DB, product *entity.Product) error {
	return tx.Debug().Updates(product).Error
}

func (r *ProductMySQL) GetTotalProductsByPartnerId(id uuid.UUID) (int64, error) {
	var total int64
	err := r.db.Debug().Model(&entity.Product{}).Where("partner_id = ?", id).Count(&total).Error
	return total, err
}

func (r *ProductMySQL) Filter(products *[]entity.Product, param dto.ProductParam, pagination dto.PaginationRequest) error {
	now := time.Now()

	query := r.db.Debug().Preload("Ratings").Where("partner_id = ? AND stock > 0 AND pickup_time >= ?", param.PartnerId, now)

	if param.Name != "" {
		query = query.Where("name LIKE ?", "%"+param.Name+"%")
	}

	query.Limit(pagination.Limit).Offset(pagination.Offset).Find(&products)

	return query.Error
}
