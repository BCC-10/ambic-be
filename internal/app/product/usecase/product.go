package usecase

import (
	partnerRepo "ambic/internal/app/partner/repository"
	"ambic/internal/app/product/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/helper"
	"ambic/internal/infra/maps"
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
	FilterProducts(request dto.FilterProductRequest) (*[]dto.GetProductResponse, *dto.PaginationResponse, *res.Err)
	ShowProduct(productId uuid.UUID) (dto.GetProductResponse, *res.Err)
	CreateProduct(userId uuid.UUID, request dto.CreateProductRequest) *res.Err
	UpdateProduct(productId uuid.UUID, partnerId uuid.UUID, req dto.UpdateProductRequest) *res.Err
	DeleteProduct(productId uuid.UUID, partnerId uuid.UUID) *res.Err
}

type ProductUsecase struct {
	env               *env.Env
	db                *gorm.DB
	ProductRepository repository.ProductMySQLItf
	PartnerRepository partnerRepo.PartnerMySQLItf
	Supabase          supabase.SupabaseIf
	Maps              maps.MapsIf
	helper            helper.HelperIf
}

func NewProductUsecase(env *env.Env, db *gorm.DB, productRepository repository.ProductMySQLItf, partnerRepository partnerRepo.PartnerMySQLItf, supabase supabase.SupabaseIf, helper helper.HelperIf, maps maps.MapsIf) ProductUsecaseItf {
	return &ProductUsecase{
		env:               env,
		db:                db,
		ProductRepository: productRepository,
		PartnerRepository: partnerRepository,
		Supabase:          supabase,
		helper:            helper,
		Maps:              maps,
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
		tx.Rollback()
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.ProductNotExists)
		}

		return res.ErrInternalServer()
	}

	if productDB.PartnerID != partnerId {
		tx.Rollback()
		return res.ErrForbidden(res.RatingNotBelongToPartner)
	}

	pickupTime, err := time.Parse("2006-01-02 15:04:05", req.PickupTime)
	if err != nil {
		tx.Rollback()
		return res.ErrBadRequest(res.InvalidDateTime)
	}

	endPickupTime, err := time.Parse("2006-01-02 15:04:05", req.PickupTime)
	if err != nil {
		tx.Rollback()
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
			tx.Rollback()
			return err
		}

		src, err := req.Photo.Open()
		if err != nil {
			tx.Rollback()
			return res.ErrInternalServer()
		}

		defer src.Close()

		bucket := u.env.SupabaseBucket
		path := "products/" + uuid.NewString() + filepath.Ext(req.Photo.Filename)
		contentType := req.Photo.Header.Get("Content-Type")

		publicURL, err := u.Supabase.UploadFile(bucket, path, contentType, src)
		if err != nil {
			tx.Rollback()
			return res.ErrInternalServer()
		}

		product.PhotoURL = publicURL

		if productDB.PhotoURL != "" {
			oldPhotoURL := productDB.PhotoURL
			index := strings.Index(oldPhotoURL, bucket)
			oldPhotoPath := oldPhotoURL[index+len(bucket+"/"):]

			if err := u.Supabase.DeleteFile(bucket, oldPhotoPath); err != nil {
				tx.Rollback()
				return res.ErrInternalServer()
			}
		}
	}

	if err := u.ProductRepository.Update(tx, product); err != nil {
		tx.Rollback()
		if mysql.CheckError(err, mysql.ErrDuplicateEntry) {
			return res.ErrBadRequest(res.ProductAlreadyExists)
		}

		return res.ErrInternalServer()
	}

	tx.Commit()

	return nil
}

func (u ProductUsecase) DeleteProduct(productId uuid.UUID, partnerId uuid.UUID) *res.Err {
	tx := u.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	productDB := new(entity.Product)
	if err := u.ProductRepository.Show(productDB, dto.ProductParam{ID: productId}); err != nil {
		tx.Rollback()
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.ProductNotExists)
		}

		return res.ErrInternalServer()
	}

	if productDB.PartnerID != partnerId {
		tx.Rollback()
		return res.ErrForbidden(res.RatingNotBelongToPartner)
	}

	if err := u.ProductRepository.Delete(tx, productDB); err != nil {
		tx.Rollback()
		return res.ErrInternalServer()
	}

	if productDB.PhotoURL != "" {
		bucket := u.env.SupabaseBucket
		index := strings.Index(productDB.PhotoURL, bucket)
		path := productDB.PhotoURL[index+len(bucket+"/"):]

		if err := u.Supabase.DeleteFile(bucket, path); err != nil {
			tx.Rollback()
			return res.ErrInternalServer()
		}
	}

	tx.Commit()

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

	return product.ParseDTOGet(nil), nil
}

func (u ProductUsecase) FilterProducts(req dto.FilterProductRequest) (*[]dto.GetProductResponse, *dto.PaginationResponse, *res.Err) {
	resp := new([]dto.GetProductResponse)

	partners := new([]entity.Partner)
	if err := u.PartnerRepository.Get(partners, dto.PartnerParam{IsVerified: true}); err != nil {
		return resp, nil, res.ErrInternalServer()
	}

	partnerDistanceMap := make(map[uuid.UUID]float64)
	var withinRadiusPartnerIds []uuid.UUID

	for _, partner := range *partners {
		origin := dto.Location{Lat: req.Lat, Long: req.Long}
		destination := dto.Location{Lat: partner.Latitude, Long: partner.Longitude}

		distance, err := u.Maps.GetDistance(origin, destination)
		if err != nil {
			return nil, nil, res.ErrInternalServer()
		}

		if distance == nil {
			continue
		}

		if float64(*distance) <= req.Radius {
			withinRadiusPartnerIds = append(withinRadiusPartnerIds, partner.ID)
			partnerDistanceMap[partner.ID] = float64(*distance)
		}
	}

	if len(withinRadiusPartnerIds) == 0 {
		return resp, nil, nil
	}

	products := new([]entity.Product)
	param := dto.ProductParam{
		PartnerIds: withinRadiusPartnerIds,
		Name:       req.Name,
	}

	pagination := u.helper.CreatePagination(dto.PaginationRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
	})

	totalProducts, err := u.ProductRepository.Filter(products, param, pagination)
	if err != nil {
		return nil, nil, res.ErrInternalServer()
	}

	if len(*products) == 0 {
		return resp, nil, nil
	}

	var response []dto.GetProductResponse
	for _, product := range *products {
		distance := partnerDistanceMap[product.PartnerID]
		response = append(response, product.ParseDTOGet(&distance))
	}

	pg := u.helper.CalculatePagination(pagination, totalProducts)

	return &response, &pg, nil
}
