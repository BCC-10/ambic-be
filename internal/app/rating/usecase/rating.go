package usecase

import (
	productRepo "ambic/internal/app/product/repository"
	"ambic/internal/app/rating/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/helper"
	"ambic/internal/infra/mysql"
	res "ambic/internal/infra/response"
	"ambic/internal/infra/supabase"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type RatingUsecaseItf interface {
	Get() (*[]dto.GetRatingResponse, *res.Err)
	Show(ratingId uuid.UUID) (dto.GetRatingResponse, *res.Err)
	Create(userId uuid.UUID, request dto.CreateRatingRequest) *res.Err
	Update(userId uuid.UUID, ratingId uuid.UUID, request dto.UpdateRatingRequest) *res.Err
	Delete(userId uuid.UUID, ratingId uuid.UUID) *res.Err
}

type RatingUsecase struct {
	env               *env.Env
	RatingRepository  repository.RatingMySQLItf
	ProductRepository productRepo.ProductMySQLItf
	Supabase          supabase.SupabaseIf
	helper            helper.HelperIf
}

func NewRatingUsecase(env *env.Env, ratingRepository repository.RatingMySQLItf, productRepository productRepo.ProductMySQLItf, supabase supabase.SupabaseIf, helper helper.HelperIf) RatingUsecaseItf {
	return &RatingUsecase{
		env:               env,
		RatingRepository:  ratingRepository,
		ProductRepository: productRepository,
		Supabase:          supabase,
		helper:            helper,
	}
}

func (u *RatingUsecase) Get() (*[]dto.GetRatingResponse, *res.Err) {
	ratings := new([]entity.Rating)
	if err := u.RatingRepository.Get(ratings); err != nil {
		return nil, res.ErrInternalServer()
	}

	resp := make([]dto.GetRatingResponse, len(*ratings))
	for i, rating := range *ratings {
		resp[i] = rating.ParseDTOGet()
	}

	return &resp, nil
}

func (u *RatingUsecase) Show(ratingId uuid.UUID) (dto.GetRatingResponse, *res.Err) {
	rating := new(entity.Rating)
	if err := u.RatingRepository.Show(rating, dto.RatingParam{ID: ratingId}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return dto.GetRatingResponse{}, res.ErrNotFound(res.RatingNotFound)
		}

		return dto.GetRatingResponse{}, res.ErrInternalServer()
	}

	return rating.ParseDTOGet(), nil
}

func (u *RatingUsecase) Create(userId uuid.UUID, request dto.CreateRatingRequest) *res.Err {
	productId, err := uuid.Parse(request.ProductID)
	if err != nil {
		return res.ErrBadRequest(res.InvalidUUID)
	}

	product := new(entity.Product)
	if err := u.ProductRepository.Show(product, dto.ProductParam{ID: productId}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.ProductNotExists)
		}

		return res.ErrInternalServer()
	}

	ratingDB := new(entity.Rating)
	if err := u.RatingRepository.Show(ratingDB, dto.RatingParam{UserID: userId, ProductID: productId}); err == nil {
		return res.ErrForbidden(res.UserAlreadyRated)
	}

	rating := &entity.Rating{
		UserID:    userId,
		ProductID: productId,
		Star:      request.Star,
		Feedback:  request.Feedback,
	}

	if request.Photo != nil {
		if err := u.helper.ValidateImage(request.Photo); err != nil {
			return err
		}

		src, err := request.Photo.Open()
		if err != nil {
			return res.ErrInternalServer()
		}

		defer src.Close()

		bucket := u.env.SupabaseBucket
		filepath := "ratings/" + uuid.NewString()
		contentType := request.Photo.Header.Get("Content-Type")

		photoURL, err := u.Supabase.UploadFile(bucket, filepath, contentType, src)

		rating.PhotoURL = photoURL
	}

	if err := u.RatingRepository.Create(rating); err != nil {
		return res.ErrInternalServer()
	}

	return nil
}

func (u *RatingUsecase) Update(userId uuid.UUID, ratingId uuid.UUID, request dto.UpdateRatingRequest) *res.Err {
	ratingDB := new(entity.Rating)
	if err := u.RatingRepository.Show(ratingDB, dto.RatingParam{ID: ratingId}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.RatingNotFound)
		}

		return res.ErrInternalServer()
	}

	if ratingDB.UserID != userId {
		return res.ErrForbidden(res.RatingNotBelongToUser)
	}

	rating := &entity.Rating{
		ID:       ratingId,
		Star:     request.Star,
		Feedback: request.Feedback,
	}

	if request.Photo != nil {
		if err := u.helper.ValidateImage(request.Photo); err != nil {
			return err
		}

		src, err := request.Photo.Open()
		if err != nil {
			return res.ErrInternalServer()
		}

		defer src.Close()

		filepath := "ratings/" + uuid.NewString()
		contentType := request.Photo.Header.Get("Content-Type")

		photoURL, err := u.Supabase.UploadFile(u.env.SupabaseBucket, filepath, contentType, src)

		rating.PhotoURL = photoURL
	}

	if err := u.RatingRepository.Update(rating); err != nil {
		return res.ErrInternalServer()
	}

	if ratingDB.PhotoURL != "" {
		oldPhotoURL := ratingDB.PhotoURL
		index := strings.Index(oldPhotoURL, u.env.SupabaseBucket)
		oldPhotoPath := oldPhotoURL[index+len(u.env.SupabaseBucket+"/"):]

		if err := u.Supabase.DeleteFile(u.env.SupabaseBucket, oldPhotoPath); err != nil {
			return res.ErrInternalServer()
		}
	}

	return nil
}

func (u *RatingUsecase) Delete(userId uuid.UUID, ratingId uuid.UUID) *res.Err {
	ratingDB := new(entity.Rating)
	if err := u.RatingRepository.Show(ratingDB, dto.RatingParam{ID: ratingId}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.RatingNotExists)
		}

		return res.ErrInternalServer()
	}

	if ratingDB.UserID != userId {
		return res.ErrForbidden(res.RatingNotBelongToUser)
	}

	if err := u.RatingRepository.Delete(ratingDB); err != nil {
		return res.ErrInternalServer()
	}

	if ratingDB.PhotoURL != "" {
		bucket := u.env.SupabaseBucket
		index := strings.Index(ratingDB.PhotoURL, bucket)
		path := ratingDB.PhotoURL[index+len(bucket+"/"):]

		if err := u.Supabase.DeleteFile(bucket, path); err != nil {
			return res.ErrInternalServer()
		}
	}

	return nil
}
