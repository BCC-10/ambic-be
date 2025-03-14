package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RatingMySQLItf interface {
	Get(rating *[]entity.Rating, param dto.RatingParam, pagination dto.PaginationRequest) (int64, error)
	Show(rating *entity.Rating, param dto.RatingParam) error
	Create(rating *entity.Rating) error
	Update(tx *gorm.DB, rating *entity.Rating) error
	Delete(tx *gorm.DB, rating *entity.Rating) error
	GetTotalRatingsByPartnerId(partnerId uuid.UUID) (int64, error)
}

type RatingMySQL struct {
	db *gorm.DB
}

func NewRatingMySQL(db *gorm.DB) RatingMySQLItf {
	return &RatingMySQL{db}
}

func (r *RatingMySQL) Get(rating *[]entity.Rating, param dto.RatingParam, pagination dto.PaginationRequest) (int64, error) {
	query := r.db.Debug().Preload("User")

	var count int64
	if err := query.Model(&rating).Count(&count).Error; err != nil {
		return 0, err
	}

	query.Limit(pagination.Limit).Offset(pagination.Offset).Order("created_at desc").Find(rating, param)

	return count, query.Error
}

func (r *RatingMySQL) Show(rating *entity.Rating, param dto.RatingParam) error {
	return r.db.Debug().Preload("User").First(rating, param).Error
}

func (r *RatingMySQL) Create(rating *entity.Rating) error {
	return r.db.Debug().Create(rating).Error
}

func (r *RatingMySQL) Update(tx *gorm.DB, rating *entity.Rating) error {
	return tx.Debug().Updates(rating).Error
}

func (r *RatingMySQL) Delete(tx *gorm.DB, rating *entity.Rating) error {
	return tx.Debug().Delete(rating).Error
}

func (r *RatingMySQL) GetTotalRatingsByPartnerId(partnerId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Table("ratings").
		Joins("JOIN products ON ratings.product_id = products.id").
		Where("products.partner_id = ?", partnerId).
		Count(&count).Error

	return count, err
}
