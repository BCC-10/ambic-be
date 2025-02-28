package usecase

type ProductUsecaseItf interface {}

type ProductUsecase struct {}

func NewProductUsecase() ProductUsecaseItf {
    return &ProductUsecase{}
}
