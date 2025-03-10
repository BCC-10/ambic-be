package usecase

import (
	"ambic/internal/domain/dto"
	"ambic/internal/infra/maps"
	res "ambic/internal/infra/response"
)

type LocationUsecaseItf interface {
	AutocompleteLocation(req dto.LocationRequest) ([]dto.LocationResponse, *res.Err)
	ShowLocation(id string) (dto.PlaceDetails, *res.Err)
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

func (u *LocationUsecase) ShowLocation(id string) (dto.PlaceDetails, *res.Err) {
	location, err := u.Maps.GetPlaceDetails(id)
	if err != nil {
		return dto.PlaceDetails{}, res.ErrNotFound()
	}

	return location, nil
}
