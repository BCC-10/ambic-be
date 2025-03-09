package usecase

import (
	"ambic/internal/app/product/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/helper"
	"ambic/internal/infra/mysql"
	res "ambic/internal/infra/response"
	"ambic/internal/infra/supabase"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
	"time"
)

type ProductUsecaseItf interface {
	ShowProduct(productId uuid.UUID) (dto.GetProductResponse, *res.Err)
	CreateProduct(userId uuid.UUID, request dto.CreateProductRequest) *res.Err
	UpdateProduct(productId uuid.UUID, partnerId uuid.UUID, req dto.UpdateProductRequest) *res.Err
	DeleteProduct(productId uuid.UUID, partnerId uuid.UUID) *res.Err
}

type ProductUsecase struct {
	env               *env.Env
	db                *gorm.DB
	ProductRepository repository.ProductMySQLItf
	Supabase          supabase.SupabaseIf
	helper            helper.HelperIf
}

func NewProductUsecase(env *env.Env, db *gorm.DB, productRepository repository.ProductMySQLItf, supabase supabase.SupabaseIf, helper helper.HelperIf) ProductUsecaseItf {
	return &ProductUsecase{
		env:               env,
		db:                db,
		ProductRepository: productRepository,
		Supabase:          supabase,
		helper:            helper,
	}
}

func (u ProductUsecase) CreateProduct(partnerId uuid.UUID, req dto.CreateProductRequest) *res.Err {
	if err := u.helper.ValidateImage(req.Photo); err != nil {
		return err
	}

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

	pickupTime, err := time.Parse("2006-01-02 15:04:05", req.PickupTime)
	if err != nil {
		return res.ErrBadRequest(res.InvalidDateTime)
	}

	endPickupTime, err := time.Parse("2006-01-02 15:04:05", req.EndPickupTime)
	if err != nil {
		return res.ErrBadRequest(res.InvalidDateTime)
	}

	product := &entity.Product{
		PartnerID:     partnerId,
		Name:          req.Name,
		Description:   req.Description,
		InitialPrice:  req.InitialPrice,
		FinalPrice:    req.FinalPrice,
		Stock:         uint(req.Stock),
		PickupTime:    pickupTime,
		EndPickupTime: endPickupTime,
		PhotoURL:      publicURL,
	}

	if err = u.ProductRepository.Create(product); err != nil {
		if mysql.CheckError(err, mysql.ErrDuplicateEntry) {
			return res.ErrBadRequest(res.ProductAlreadyExists)
		}
		return res.ErrInternalServer()
	}

	return nil
}

func (u ProductUsecase) UpdateProduct(productId uuid.UUID, partnerId uuid.UUID, req dto.UpdateProductRequest) *res.Err {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	productDB := new(entity.Product)
	if err := u.ProductRepository.Show(productDB, dto.ProductParam{ID: productId}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.ProductNotExists)
		}

		return res.ErrInternalServer()
	}

	if productDB.PartnerID != partnerId {
		return res.ErrForbidden(res.RatingNotBelongToPartner)
	}

	pickupTime, err := time.Parse("2006-01-02 15:04:05", req.PickupTime)
	if err != nil {
		return res.ErrBadRequest(res.InvalidDateTime)
	}

	endPickupTime, err := time.Parse("2006-01-02 15:04:05", req.PickupTime)
	if err != nil {
		return res.ErrBadRequest(res.InvalidDateTime)
	}

	product := &entity.Product{
		ID:            productId,
		Name:          req.Name,
		Description:   req.Description,
		InitialPrice:  req.InitialPrice,
		FinalPrice:    req.FinalPrice,
		Stock:         uint(req.Stock),
		PickupTime:    pickupTime,
		EndPickupTime: endPickupTime,
	}

	if req.Photo != nil {
		if err := u.helper.ValidateImage(req.Photo); err != nil {
			return err
		}

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

	if err := u.ProductRepository.Update(tx, product); err != nil {
		if mysql.CheckError(err, mysql.ErrDuplicateEntry) {
			return res.ErrBadRequest(res.ProductAlreadyExists)
		}

		return res.ErrInternalServer()
	}

	tx.Commit()

	return nil
}

func (u ProductUsecase) DeleteProduct(productId uuid.UUID, partnerId uuid.UUID) *res.Err {
	productDB := new(entity.Product)
	if err := u.ProductRepository.Show(productDB, dto.ProductParam{ID: productId}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.ProductNotExists)
		}

		return res.ErrInternalServer()
	}

	if productDB.PartnerID != partnerId {
		return res.ErrForbidden(res.RatingNotBelongToPartner)
	}

	if err := u.ProductRepository.Delete(productDB); err != nil {
		return res.ErrInternalServer()
	}

	if productDB.PhotoURL != "" {
		bucket := u.env.SupabaseBucket
		index := strings.Index(productDB.PhotoURL, bucket)
		path := productDB.PhotoURL[index+len(bucket+"/"):]

		if err := u.Supabase.DeleteFile(bucket, path); err != nil {
			return res.ErrInternalServer()
		}
	}

	return nil
}

func (u ProductUsecase) ShowProduct(productId uuid.UUID) (dto.GetProductResponse, *res.Err) {
	product := new(entity.Product)
	if err := u.ProductRepository.Show(product, dto.ProductParam{ID: productId}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return dto.GetProductResponse{}, res.ErrNotFound(res.ProductNotExists)
		}

		return dto.GetProductResponse{}, res.ErrInternalServer()
	}

	return product.ParseDTOGet(), nil
}
