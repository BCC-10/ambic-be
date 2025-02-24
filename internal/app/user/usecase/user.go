package usecase

type UserUsecaseItf interface {}

type UserUsecase struct {}

func NewUserUsecase() UserUsecaseItf {
    return &UserUsecase{}
}
