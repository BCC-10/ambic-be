package repository

import(
    "gorm.io/gorm"
)

type ProductMySQLItf interface {}

type ProductMySQL struct {
    db *gorm.DB
}

func NewProductMySQL(db *gorm.DB) ProductMySQLItf {
    return &ProductMySQL{db}
}
