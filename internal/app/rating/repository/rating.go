package repository

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"gorm.io/gorm"
)

type RatingMySQLItf interface {
	Get(rating *entity.Rating) error
	Show(rating *entity.Rating, param dto.RatingParam) error
	Create(rating *entity.Rating) error
	Update(rating *entity.Rating) error
	Delete(rating *entity.Rating) error
}

type RatingMySQL struct {
	db *gorm.DB
}

func NewRatingMySQL(db *gorm.DB) RatingMySQLItf {
	return &RatingMySQL{db}
}

func (r *RatingMySQL) Get(rating *entity.Rating) error {
	return r.db.Debug().Find(rating).Error
}

func (r *RatingMySQL) Show(rating *entity.Rating, param dto.RatingParam) error {
	return r.db.Debug().First(rating, param).Error
}

func (r *RatingMySQL) Create(rating *entity.Rating) error {
	return r.db.Debug().Create(rating).Error
}

func (r *RatingMySQL) Update(rating *entity.Rating) error {
	return r.db.Debug().Updates(rating).Error
}

func (r *RatingMySQL) Delete(rating *entity.Rating) error {
	return r.db.Debug().Delete(rating).Error
}
