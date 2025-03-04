package usecase

import (
	"ambic/internal/app/product/repository"
	userRepo "ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/mysql"
	res "ambic/internal/infra/response"
	"ambic/internal/infra/supabase"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
)

type ProductUsecaseItf interface {
	CreateProduct(userId uuid.UUID, request dto.CreateProductRequest) *res.Err
	UpdateProduct(productId uuid.UUID, partnerId uuid.UUID, req dto.UpdateProductRequest) *res.Err
}

type ProductUsecase struct {
	env               *env.Env
	UserRepository    userRepo.UserMySQLItf
	ProductRepository repository.ProductMySQLItf
	Supabase          supabase.SupabaseIf
}

func NewProductUsecase(env *env.Env, productRepository repository.ProductMySQLItf, userRepository userRepo.UserMySQLItf, supabase supabase.SupabaseIf) ProductUsecaseItf {
	return &ProductUsecase{
		env:               env,
		UserRepository:    userRepository,
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
	if err := u.UserRepository.Get(user, dto.UserParam{Id: userId}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.UserNotExists)
		}

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

	if err = u.ProductRepository.Create(product); err != nil {
		return res.ErrInternalServer()
	}

	return nil
}

func (u ProductUsecase) UpdateProduct(productId uuid.UUID, partnerId uuid.UUID, req dto.UpdateProductRequest) *res.Err {
	productDB := new(entity.Product)
	if err := u.ProductRepository.Show(productDB, dto.ProductParam{Id: productId}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.ProductNotExists)
		}

		return res.ErrInternalServer()
	}

	if productDB.PartnerID != partnerId {
		return res.ErrForbidden(res.ProductNotBelongToPartner)
	}

	product := &entity.Product{
		ID:           productId,
		Name:         req.Name,
		Description:  req.Description,
		InitialPrice: req.InitialPrice,
		FinalPrice:   req.FinalPrice,
		Stock:        req.Stock,
		PickupTime:   req.PickupTime,
	}

	if req.Photo != nil {
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

		product.PhotoURL = publicURL

		if productDB.PhotoURL != "" {
			oldPhotoURL := productDB.PhotoURL
			index := strings.Index(oldPhotoURL, bucket)
			oldPhotoPath := oldPhotoURL[index+len(bucket+"/"):]

			if err := u.Supabase.DeleteFile(bucket, oldPhotoPath); err != nil {
				return res.ErrInternalServer()
			}
		}
	}

	if err := u.ProductRepository.Update(product); err != nil {
		return res.ErrInternalServer()
	}

	return nil
}
