package usecase

import (
	"ambic/internal/app/partner/repository"
	userRepo "ambic/internal/app/user/repository"
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
)

type PartnerUsecaseItf interface {
	RegisterPartner(id uuid.UUID, data dto.RegisterPartnerRequest) *res.Err
	VerifyPartner(request dto.VerifyPartnerRequest) *res.Err
	GetProducts(id uuid.UUID, query dto.GetPartnerProductsQuery) ([]dto.GetProductResponse, *res.Err)
}

type PartnerUsecase struct {
	env               *env.Env
	PartnerRepository repository.PartnerMySQLItf
	UserRepository    userRepo.UserMySQLItf
	Maps              maps.MapsIf
	Supabase          supabase.SupabaseIf
	helper            helper.HelperIf
}

func NewPartnerUsecase(env *env.Env, partnerRepository repository.PartnerMySQLItf, userRepository userRepo.UserMySQLItf, supabase supabase.SupabaseIf, helper helper.HelperIf, maps maps.MapsIf) PartnerUsecaseItf {
	return &PartnerUsecase{
		env:               env,
		PartnerRepository: partnerRepository,
		Maps:              maps,
		UserRepository:    userRepository,
		Supabase:          supabase,
		helper:            helper,
	}
}

func (u *PartnerUsecase) RegisterPartner(id uuid.UUID, data dto.RegisterPartnerRequest) *res.Err {
	user := new(entity.User)
	if err := u.UserRepository.Show(user, dto.UserParam{Id: id}); err != nil {
		return res.ErrInternalServer()
	}

	if user.Name == "" || user.Phone == "" || user.Address == "" || user.Gender == nil || user.BornDate.IsZero() {
		return res.ErrForbidden(res.ProfileNotFilledCompletely)
	}

	if data.Instagram[0] == '@' {
		data.Instagram = data.Instagram[1:]
	}

	partner := entity.Partner{
		UserID:    id,
		Name:      data.Name,
		Type:      data.Type,
		Address:   data.Address,
		City:      data.City,
		Instagram: data.Instagram,
		Longitude: data.Longitude,
		Latitude:  data.Latitude,
	}

	if data.Photo != nil {
		if err := u.helper.ValidateImage(data.Photo); err != nil {
			return err
		}

		src, err := data.Photo.Open()
		if err != nil {
			return res.ErrInternalServer()
		}

		defer src.Close()

		bucket := u.env.SupabaseBucket
		path := "partners/" + uuid.NewString() + filepath.Ext(data.Photo.Filename)
		contentType := data.Photo.Header.Get("Content-Type")

		photoURL, err := u.Supabase.UploadFile(bucket, path, contentType, src)
		if err != nil {
			return res.ErrInternalServer()
		}

		partner.PhotoURL = photoURL
	}

	if err := u.PartnerRepository.Create(&partner); err != nil {
		return res.ErrInternalServer()
	}

	return nil
}

func (u *PartnerUsecase) VerifyPartner(data dto.VerifyPartnerRequest) *res.Err {
	user := new(entity.User)
	if err := u.UserRepository.Show(user, dto.UserParam{Email: data.Email}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.UserNotExists)
		}

		return res.ErrInternalServer()
	}

	if user.Partner.ID == uuid.Nil {
		return res.ErrNotFound(res.PartnerNotExists)
	}

	if user.Partner.IsVerified {
		return res.ErrForbidden(res.PartnerVerified)
	}

	if data.Token != u.env.PartnerVerificationToken {
		return res.ErrForbidden(res.InvalidToken)
	}

	user.Partner.IsVerified = true

	if err := u.PartnerRepository.Update(&user.Partner); err != nil {
		return res.ErrInternalServer()
	}

	return nil
}

func (u *PartnerUsecase) GetProducts(id uuid.UUID, query dto.GetPartnerProductsQuery) ([]dto.GetProductResponse, *res.Err) {
	if query.Limit < 1 {
		query.Limit = 1
	}

	if query.Page < 1 {
		query.Page = 1
	}

	limit := query.Limit
	offset := (query.Page - 1) * query.Limit

	partner := new(entity.Partner)
	if err := u.PartnerRepository.GetProducts(partner, dto.PartnerParam{ID: id}, limit, offset); err != nil {
		return nil, res.ErrInternalServer()
	}

	products := make([]dto.GetProductResponse, 0)
	for _, product := range partner.Products {
		products = append(products, dto.GetProductResponse{
			ID:           product.ID.String(),
			Name:         product.Name,
			Description:  product.Description,
			InitialPrice: product.InitialPrice,
			FinalPrice:   product.FinalPrice,
			Stock:        product.Stock,
			PickupTime:   product.PickupTime,
			PhotoURL:     product.PhotoURL,
		})
	}

	return products, nil
}
