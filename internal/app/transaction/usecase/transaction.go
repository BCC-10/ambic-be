package usecase

type TransactionUsecaseItf interface{}

type TransactionUsecase struct{}

func NewTransactionUsecase() TransactionUsecaseItf {
	return &TransactionUsecase{}
}
