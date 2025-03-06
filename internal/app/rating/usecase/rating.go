package usecase

import (
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
	Create(userId uuid.UUID, request dto.CreateRatingRequest) *res.Err
	Update(userId uuid.UUID, ratingId uuid.UUID, request dto.UpdateRatingRequest) *res.Err
	Delete(userId uuid.UUID, ratingId uuid.UUID) *res.Err
}

type RatingUsecase struct {
	env              *env.Env
	RatingRepository repository.RatingMySQLItf
	Supabase         supabase.SupabaseIf
	helper           helper.HelperIf
}

func NewRatingUsecase(env *env.Env, ratingRepository repository.RatingMySQLItf, supabase supabase.SupabaseIf, helper helper.HelperIf) RatingUsecaseItf {
	return &RatingUsecase{
		env:              env,
		RatingRepository: ratingRepository,
		Supabase:         supabase,
		helper:           helper,
	}
}

func (u *RatingUsecase) Get() (*[]dto.GetRatingResponse, *res.Err) {
	ratings := new([]entity.Rating)
	if err := u.RatingRepository.Get(ratings); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return nil, res.ErrNotFound(res.RatingEmpty)
		}
	}

	resp := make([]dto.GetRatingResponse, len(*ratings))
	for i, rating := range *ratings {
		resp[i] = rating.ParseDTOGet()
	}

	return &resp, nil
}

func (u *RatingUsecase) Create(userId uuid.UUID, request dto.CreateRatingRequest) *res.Err {
	productId, _ := uuid.Parse(request.ProductID)

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
			return res.ErrNotFound(res.ProductNotExists)
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
