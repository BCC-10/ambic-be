package entity

type BusinessType struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255);not null"`
}
