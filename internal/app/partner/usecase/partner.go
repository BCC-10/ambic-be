package usecase

import (
	repository "ambic/internal/app/partner/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/maps"
	res "ambic/internal/infra/response"
	"github.com/google/uuid"
)

type PartnerUsecaseItf interface {
	RegisterPartner(id uuid.UUID, data dto.RegisterPartnerRequest) *res.Err
}

type PartnerUsecase struct {
	env               *env.Env
	PartnerRepository repository.PartnerMySQLItf
	Maps              maps.MapsIf
}

func NewPartnerUsecase(env *env.Env, partnerRepository repository.PartnerMySQLItf, maps maps.MapsIf) PartnerUsecaseItf {
	return &PartnerUsecase{
		env:               env,
		PartnerRepository: partnerRepository,
		Maps:              maps,
	}
}

func (u *PartnerUsecase) RegisterPartner(id uuid.UUID, data dto.RegisterPartnerRequest) *res.Err {
	partner := entity.Partner{
		ID:        uuid.New(),
		UserID:    id,
		Name:      data.Name,
		Type:      data.Type,
		Address:   data.Address,
		City:      data.City,
		Longitude: data.Longitude,
		Latitude:  data.Latitude,
	}

	err := u.PartnerRepository.Create(&partner)
	if err != nil {
		return res.ErrInternalServer()
	}

	return nil
}
