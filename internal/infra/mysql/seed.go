package mysql

import (
	BusinessTypeRepo "ambic/internal/app/business_type/repository"
	PaymentMethodRepo "ambic/internal/app/payment_method/repository"
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

	paymentMethodRepository := PaymentMethodRepo.NewPaymentMethodMySQL(db)

	paymentMethodNames := []string{
		"credit_card",
		"gopay",
		"qris",
		"shopeepay",
		"bank_transfer",
		"e_channel",
		"cstore",
		"akulaku",
		"others",
	}

	for _, name := range paymentMethodNames {
		paymentMethod := &entity.PaymentMethod{
			Name: name,
		}

		if err := paymentMethodRepository.Create(paymentMethod); err != nil {
			fmt.Printf("Failed to add payment method: %s, error: %v\n", name, err)
		} else {
			fmt.Printf("Success adding payment method: %s\n", name)
		}
	}
}
