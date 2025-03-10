package usecase

import (
	"ambic/internal/app/payment/repository"
	ProductRepo "ambic/internal/app/product/repository"
	transactionRepo "ambic/internal/app/transaction/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/infra/mysql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"time"

	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
)

type PaymentUsecaseItf interface {
	ProcessPayment(req *dto.NotificationPayment) *res.Err
}

type PaymentUsecase struct {
	env                   *env.Env
	db                    *gorm.DB
	PaymentRepository     repository.PaymentMySQLItf
	TransactionRepository transactionRepo.TransactionMySQLItf
	ProductRepository     ProductRepo.ProductMySQLItf
}

func NewPaymentUsecase(env *env.Env, db *gorm.DB, paymentRepository repository.PaymentMySQLItf, transactionRepository transactionRepo.TransactionMySQLItf, productRepository ProductRepo.ProductMySQLItf) PaymentUsecaseItf {
	return &PaymentUsecase{
		env:                   env,
		db:                    db,
		PaymentRepository:     paymentRepository,
		TransactionRepository: transactionRepository,
		ProductRepository:     productRepository,
	}
}

func (u PaymentUsecase) ProcessPayment(req *dto.NotificationPayment) *res.Err {
	tx := u.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transactionStatusResp := req.TransactionStatus
	fraudStatus := req.FraudStatus

	transactionTime, _ := time.Parse("2006-01-02 15:04:05", req.TransactionTime)
	settlementTime, _ := time.Parse("2006-01-02 15:04:05", req.SettlementTime)

	transactionId, err := uuid.Parse(req.TransactionID)
	if err != nil {
		return res.ErrBadRequest()
	}

	grossAmount, _ := strconv.ParseFloat(req.GrossAmount, 32)

	payment := &entity.Payment{
		TransactionID:     transactionId,
		OrderID:           req.OrderID,
		ReferenceID:       req.ReferenceID,
		MerchantID:        req.MerchantID,
		Issuer:            req.Issuer,
		GrossAmount:       float32(grossAmount),
		Currency:          req.Currency,
		Acquirer:          req.Acquirer,
		TransactionStatus: req.TransactionStatus,
		StatusMessage:     req.StatusMessage,
		PaymentType:       req.PaymentType,
		FraudStatus:       req.FraudStatus,
		TransactionTime:   transactionTime,
		SettlementTime:    settlementTime,
	}

	transaction := new(entity.Transaction)
	transaction.ID = transactionId

	var status string

	if transactionStatusResp == "settlement" || (transactionStatusResp == "capture" && fraudStatus == "accept") {
		status = "success"
	} else if transactionStatusResp == "cancel" || transactionStatusResp == "expire" {
		status = "cancelled"
	}

	if status == "success" {
		paymentDB := new(entity.Payment)
		if err := u.PaymentRepository.Show(paymentDB, dto.PaymentParam{TransactionID: transactionId}); err != nil {
			if mysql.CheckError(err, gorm.ErrRecordNotFound) {
				if err := u.PaymentRepository.Create(tx, payment); err != nil {
					tx.Rollback()
					return res.ErrInternalServer()
				}
			} else {
				tx.Rollback()
				return res.ErrInternalServer()
			}
		}

		transaction.Status = entity.Process
		transaction.Payment = *payment

		if err := u.TransactionRepository.Update(tx, transaction); err != nil {
			tx.Rollback()
			return res.ErrInternalServer()
		}
	} else if status == "cancelled" {
		transaction.Status = entity.CancelledBySystem

		if err := u.TransactionRepository.Update(tx, transaction); err != nil {
			tx.Rollback()
			return res.ErrInternalServer()
		}

		transactionDB := new(entity.Transaction)
		if err := u.TransactionRepository.Show(transactionDB, dto.TransactionParam{ID: transactionId}); err != nil {
			tx.Rollback()
			return res.ErrInternalServer()
		}

		for _, details := range transactionDB.TransactionDetails {
			product := new(entity.Product)

			if err := u.ProductRepository.Show(product, dto.ProductParam{ID: details.ProductID}); err != nil {
				tx.Rollback()
				return res.ErrInternalServer()
			}

			product = &entity.Product{
				ID:    details.ProductID,
				Stock: product.Stock + details.Qty,
			}

			if err := u.ProductRepository.Update(tx, product); err != nil {
				tx.Rollback()
				return res.ErrInternalServer()
			}
		}
	}

	tx.Commit()

	return nil
}
