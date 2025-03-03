package usecase

import (
	"ambic/internal/app/product/repository"
	userRepo "ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
	"ambic/internal/infra/supabase"
	"github.com/google/uuid"
	"path/filepath"
)

type ProductUsecaseItf interface {
	CreateProduct(userId uuid.UUID, request dto.CreateProductRequest) *res.Err
}

type ProductUsecase struct {
	env               *env.Env
	UserRepository    userRepo.UserMySQLItf
	ProductRepository repository.ProductMySQLItf
	Supabase          supabase.SupabaseIf
}

func NewProductUsecase(env *env.Env, productRepository repository.ProductMySQLItf, userRepostiory userRepo.UserMySQLItf, supabase supabase.SupabaseIf) ProductUsecaseItf {
	return &ProductUsecase{
		env:               env,
		UserRepository:    userRepostiory,
		ProductRepository: productRepository,
		Supabase:          supabase,
	}
}

func (u ProductUsecase) CreateProduct(userId uuid.UUID, req dto.CreateProductRequest) *res.Err {
	src, err := req.Photo.Open()
	if err != nil {
		return res.ErrInternalServer()
	}

	defer src.Close()

	bucket := u.env.SupabaseBucket
	path := "products/" + uuid.NewString() + filepath.Ext(req.Photo.Filename)
	contentType := req.Photo.Header.Get("Content-Type")

	publicURL, err := u.Supabase.UploadFile(bucket, path, contentType, src)
	if err != nil {
		return res.ErrInternalServer()
	}

	user := new(entity.User)
	err = u.UserRepository.Get(user, dto.UserParam{Id: userId})
	if err != nil {
		return res.ErrInternalServer()
	}

	partnerId := user.Partner.ID

	product := &entity.Product{
		PartnerID:    partnerId,
		Name:         req.Name,
		Description:  req.Description,
		InitialPrice: req.InitialPrice,
		FinalPrice:   req.FinalPrice,
		Stock:        req.Stock,
		PickupTime:   req.PickupTime,
		PhotoURL:     publicURL,
	}

	err = u.ProductRepository.Create(product)
	if err != nil {
		return res.ErrInternalServer()
	}

	return nil
}
