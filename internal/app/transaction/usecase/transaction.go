package usecase

import (
	productRepo "ambic/internal/app/product/repository"
	"ambic/internal/app/transaction/repository"
	userRepo "ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/helper"
	"ambic/internal/infra/midtrans"
	"ambic/internal/infra/mysql"
	res "ambic/internal/infra/response"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionUsecaseItf interface {
	GetByUserID(userId uuid.UUID) (*[]dto.GetTransactionResponse, *res.Err)
	Create(id uuid.UUID, req *dto.CreateTransactionRequest) (string, *res.Err)
}

type TransactionUsecase struct {
	db                    *gorm.DB
	env                   *env.Env
	TransactionRepository repository.TransactionMySQLItf
	ProductRepository     productRepo.ProductMySQLItf
	UserRepository        userRepo.UserMySQLItf
	helper                helper.HelperIf
	Snap                  midtrans.MidtransIf
}

func NewTransactionUsecase(env *env.Env, db *gorm.DB, transactionRepository repository.TransactionMySQLItf, productRepository productRepo.ProductMySQLItf, userRepository userRepo.UserMySQLItf, helper helper.HelperIf, snap midtrans.MidtransIf) TransactionUsecaseItf {
	return &TransactionUsecase{
		db:                    db,
		env:                   env,
		TransactionRepository: transactionRepository,
		ProductRepository:     productRepository,
		UserRepository:        userRepository,
		helper:                helper,
		Snap:                  snap,
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

func (u *TransactionUsecase) Create(userId uuid.UUID, req *dto.CreateTransactionRequest) (string, *res.Err) {
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

	if len(req.TransactionDetails) < 1 {
		return "", res.ErrBadRequest(res.MissingTransactionItems)
	}

	var items []dto.TransactionDetail

	for _, item := range req.TransactionDetails {
		if item.Qty < 1 {
			return "", res.ErrBadRequest(res.InvalidQty)
		}

		if item.ProductID == "" {
			return "", res.ErrBadRequest(res.MissingProductID)
		}

		product := new(entity.Product)

		productId, err := uuid.Parse(item.ProductID)
		if err != nil {
			tx.Rollback()
			return "", res.ErrBadRequest(res.InvalidUUID)
		}

		if err := u.ProductRepository.Show(product, dto.ProductParam{ID: productId}); err != nil {
			if mysql.CheckError(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return "", res.ErrBadRequest(fmt.Sprintf(res.ProductNotFound, item.ProductID))
			}

			tx.Rollback()
			return "", res.ErrInternalServer()
		}

		if product.Stock < uint(item.Qty) {
			return "", res.ErrBadRequest(fmt.Sprintf(res.InsufficientStock, product.Name))
		}

		transactionDetail := entity.TransactionDetail{
			TransactionID: transactionId,
			ProductID:     productId,
			Qty:           uint(item.Qty),
		}

		transaction.Total += product.FinalPrice * float32(uint(item.Qty))
		transaction.TransactionDetails = append(transaction.TransactionDetails, transactionDetail)

		items = append(items, dto.TransactionDetail{
			ProductID: item.ProductID,
			Product:   product.ParseDTOGet(),
			Qty:       uint(item.Qty),
		})

		product = &entity.Product{ID: productId, Stock: product.Stock - uint(item.Qty)}

		if err := u.ProductRepository.Update(tx, product); err != nil {
			tx.Rollback()
			return "", res.ErrInternalServer()
		}
	}

	user := new(entity.User)
	if err := u.UserRepository.Show(user, dto.UserParam{Id: userId}); err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer()
	}

	snapReq := dto.RequestSnap{
		TransactionID:      transactionId.String(),
		OrderID:            transaction.Invoice,
		Amount:             int64(transaction.Total),
		User:               user.ParseDTOGet(),
		TransactionDetails: items,
	}

	url, err := u.Snap.GeneratePaymentLink(snapReq)
	if err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer(err.Error())
	}

	if err := u.TransactionRepository.Create(tx, transaction); err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer()
	}

	tx.Commit()

	return url, nil
}
