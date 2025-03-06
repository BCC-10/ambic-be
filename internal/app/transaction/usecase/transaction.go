package usecase

import (
	"ambic/internal/app/transaction/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
	"github.com/google/uuid"
)

type TransactionUsecaseItf interface {
	GetByUserID(userId uuid.UUID) (*[]dto.GetTransactionResponse, *res.Err)
}

type TransactionUsecase struct {
	env                   *env.Env
	TransactionRepository repository.TransactionMySQLItf
}

func NewTransactionUsecase(env *env.Env, transactionRepository repository.TransactionMySQLItf) TransactionUsecaseItf {
	return &TransactionUsecase{
		env:                   env,
		TransactionRepository: transactionRepository,
	}
}

func (u *TransactionUsecase) GetByUserID(userId uuid.UUID) (*[]dto.GetTransactionResponse, *res.Err) {
	transactions := new([]entity.Transaction)

	if err := u.TransactionRepository.Get(transactions, dto.TransactionParam{UserID: userId}); err != nil {
		return nil, nil
	}

	resp := make([]dto.GetTransactionResponse, len(*transactions))
	for i, transaction := range *transactions {
		resp[i] = transaction.ParseDTOGet()
	}

	return &resp, nil
}
