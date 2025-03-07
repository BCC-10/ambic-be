package usecase

import (
	productRepo "ambic/internal/app/product/repository"
	"ambic/internal/app/transaction/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/helper"
	res "ambic/internal/infra/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionUsecaseItf interface {
	GetByUserID(userId uuid.UUID) (*[]dto.GetTransactionResponse, *res.Err)
	Create(id uuid.UUID, req *dto.CreateTransactionRequest) *res.Err
}

type TransactionUsecase struct {
	db                    *gorm.DB
	env                   *env.Env
	TransactionRepository repository.TransactionMySQLItf
	ProductRepository     productRepo.ProductMySQLItf
	helper                helper.HelperIf
}

func NewTransactionUsecase(env *env.Env, db *gorm.DB, transactionRepository repository.TransactionMySQLItf, productRepository productRepo.ProductMySQLItf, helper helper.HelperIf) TransactionUsecaseItf {
	return &TransactionUsecase{
		db:                    db,
		env:                   env,
		TransactionRepository: transactionRepository,
		ProductRepository:     productRepository,
		helper:                helper,
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

func (u *TransactionUsecase) Create(userId uuid.UUID, req *dto.CreateTransactionRequest) *res.Err {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transactionId, _ := uuid.NewV7()

	transaction := &entity.Transaction{
		ID:      transactionId,
		UserID:  userId,
		Invoice: u.helper.GenerateInvoiceNumber(),
		Total:   0,
		Status:  entity.WaitingForPayment,
		Note:    req.Note,
	}

	for _, item := range req.TransactionDetails {
		product := new(entity.Product)

		productId, err := uuid.Parse(item.ProductID)
		if err != nil {
			tx.Rollback()
			return res.ErrBadRequest()
		}

		if err := u.ProductRepository.Show(product, dto.ProductParam{ID: productId}); err != nil {
			tx.Rollback()
			return res.ErrBadRequest()
		}

		if product.Stock < uint(item.Qty) {
			return res.ErrBadRequest(res.InsufficentStock + product.Name)
		}

		transactionDetail := entity.TransactionDetail{
			TransactionID: transactionId,
			ProductID:     productId,
			Qty:           uint(item.Qty),
		}

		transaction.Total += product.FinalPrice * float32(uint(item.Qty))
		transaction.TransactionDetails = append(transaction.TransactionDetails, transactionDetail)

		product = &entity.Product{ID: productId, Stock: product.Stock - uint(item.Qty)}

		if err := u.ProductRepository.Update(tx, product); err != nil {
			tx.Rollback()
			return res.ErrInternalServer()
		}
	}

	if err := u.TransactionRepository.Create(tx, transaction); err != nil {
		tx.Rollback()
		return res.ErrInternalServer()
	}

	tx.Commit()

	return nil
}
