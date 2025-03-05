package usecase

type RatingUsecaseItf interface{}

type RatingUsecase struct{}

func NewRatingUsecase() RatingUsecaseItf {
	return &RatingUsecase{}
}
