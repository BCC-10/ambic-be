package usecase

import (
	repository "ambic/internal/app/partner/repository"
	userRepo "ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/maps"
	"ambic/internal/infra/mysql"
	res "ambic/internal/infra/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PartnerUsecaseItf interface {
	RegisterPartner(id uuid.UUID, data dto.RegisterPartnerRequest) *res.Err
	VerifyPartner(request dto.VerifyPartnerRequest) *res.Err
}

type PartnerUsecase struct {
	env               *env.Env
	PartnerRepository repository.PartnerMySQLItf
	UserRepository    userRepo.UserMySQLItf
	Maps              maps.MapsIf
}

func NewPartnerUsecase(env *env.Env, partnerRepository repository.PartnerMySQLItf, userRepository userRepo.UserMySQLItf, maps maps.MapsIf) PartnerUsecaseItf {
	return &PartnerUsecase{
		env:               env,
		PartnerRepository: partnerRepository,
		Maps:              maps,
		UserRepository:    userRepository,
	}
}

func (u *PartnerUsecase) RegisterPartner(id uuid.UUID, data dto.RegisterPartnerRequest) *res.Err {
	user := new(entity.User)
	if err := u.UserRepository.Get(user, dto.UserParam{Id: id}); err != nil {
		return res.ErrInternalServer()
	}

	if user.Name == "" || user.Phone == "" || user.Address == "" || user.Gender == nil || user.BornDate.IsZero() {
		return res.ErrForbidden(res.ProfileNotFilledCompletely)
	}

	partner := entity.Partner{
		UserID:    id,
		Name:      data.Name,
		Type:      data.Type,
		Address:   data.Address,
		City:      data.City,
		Longitude: data.Longitude,
		Latitude:  data.Latitude,
	}

	if err := u.PartnerRepository.Create(&partner); err != nil {
		return res.ErrInternalServer()
	}

	return nil
}

func (u *PartnerUsecase) VerifyPartner(data dto.VerifyPartnerRequest) *res.Err {
	user := new(entity.User)
	if err := u.UserRepository.Get(user, dto.UserParam{Email: data.Email}); err != nil {
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
