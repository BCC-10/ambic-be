package mysql

import (
	BusinessTypeRepo "ambic/internal/app/business_type/repository"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"fmt"
)

func Seed() {
	config, err := env.New()
	if err != nil {
		panic(err)
	}

	db, err := New(fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	))

	if err != nil {
		panic(err)
	}

	businessTypeRepository := BusinessTypeRepo.NewBusinessTypeMySQL(db)

	businessTypeNames := []string{
		"Restoran/Rumah Makan",
		"Kafe/Coffe Shop",
		"Bakery/Toko Roti",
		"Catering/Jasa Makanan",
		"Hotel/Resort",
		"Food Court/Street Food",
	}

	for _, name := range businessTypeNames {
		businessType := &entity.BusinessType{
			Name: name,
		}

		if err := businessTypeRepository.Create(businessType); err != nil {
			fmt.Printf("Failed to add business type: %s, error: %v\n", name, err)
		} else {
			fmt.Printf("Success adding business type: %s\n", name)
		}
	}
}
