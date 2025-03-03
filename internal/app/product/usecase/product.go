package usecase

import (
	"ambic/internal/app/product/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
)

type ProductUsecaseItf interface {
	CreateProduct(request dto.CreateProductRequest) *res.Err
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

func (u ProductUsecase) CreateProduct(req dto.CreateProductRequest) *res.Err {
	product := &entity.Product{
		PartnerID:    req.PartnerID,
		Name:         req.Name,
		Description:  req.Description,
		InitialPrice: req.InitialPrice,
		FinalPrice:   req.FinalPrice,
		Stock:        req.Stock,
		PickupTime:   req.PickupTime,
		PhotoURL:     req.Photo,
	}

	err := u.ProductRepository.Create(product)
	if err != nil {
		return res.ErrInternalServer()
	}

	return nil
}
