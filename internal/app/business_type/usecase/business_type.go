package usecase

import (
	"ambic/internal/app/business_type/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
)

type BusinessTypeUsecaseItf interface {
	Get() (*[]dto.GetBusinessTypeResponse, *res.Err)
}

type BusinessTypeUsecase struct {
	env                    *env.Env
	BusinessTypeRepository repository.BusinessTypeMySQLItf
}

func NewBusinessTypeUsecase(env *env.Env, businessTypeRepository repository.BusinessTypeMySQLItf) BusinessTypeUsecaseItf {
	return &BusinessTypeUsecase{
		env:                    env,
		BusinessTypeRepository: businessTypeRepository,
	}
}

func (u *BusinessTypeUsecase) Get() (*[]dto.GetBusinessTypeResponse, *res.Err) {
	businessTypes := new([]entity.BusinessType)
	if err := u.BusinessTypeRepository.Get(businessTypes); err != nil {
		return nil, res.ErrNotFound(res.BusinessTypeEmpty)
	}

	resp := make([]dto.GetBusinessTypeResponse, len(*businessTypes))
	for i, businessType := range *businessTypes {
		resp[i] = businessType.ParseDTOGet()
	}

	return &resp, nil
}
