package repository

import (
	"gorm.io/gorm"
)

type RatingMySQLItf interface{}

type RatingMySQL struct {
	db *gorm.DB
}

func NewRatingMySQL(db *gorm.DB) RatingMySQLItf {
	return &RatingMySQL{db}
}
