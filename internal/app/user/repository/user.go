package repository

import(
    "gorm.io/gorm"
)

type UserMySQLItf interface {}

type UserMySQL struct {
    db *gorm.DB
}

func NewUserMySQL(db *gorm.DB) UserMySQLItf {
    return &UserMySQL{db}
}
