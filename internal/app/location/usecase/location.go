package usecase

import (
	"ambic/internal/domain/dto"
	"ambic/internal/infra/maps"
	res "ambic/internal/infra/response"
)

type LocationUsecaseItf interface {
	AutocompleteLocation(req dto.LocationRequest) ([]dto.LocationResponse, *res.Err)
}

type LocationUsecase struct {
	Maps maps.MapsIf
}

func NewLocationUsecase(maps maps.MapsIf) *LocationUsecase {
	return &LocationUsecase{
		Maps: maps,
	}
}

func (u *LocationUsecase) AutocompleteLocation(req dto.LocationRequest) ([]dto.LocationResponse, *res.Err) {
	suggestions, err := u.Maps.GetAutocomplete(req)
	if err != nil {
		return nil, res.ErrInternalServer()
	}

	return suggestions, nil
}
