package usecase

import (
	"ambic/internal/app/product/repository"
	"ambic/internal/domain/env"
)

type ProductUsecaseItf interface {
}

type ProductUsecase struct {
	env               *env.Env
	ProductRepository repository.ProductMySQLItf
}

func NewProductUsecase(env *env.Env, productRepository repository.ProductMySQLItf) ProductUsecaseItf {
	return &ProductUsecase{
		env:               env,
		ProductRepository: productRepository,
	}
}
