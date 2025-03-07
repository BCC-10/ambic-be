package usecase

import (
	"ambic/internal/app/payment/repository"
	transactionRepo "ambic/internal/app/transaction/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"github.com/google/uuid"
	"strings"
	"time"

	//"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
	//"time"
)

type PaymentUsecaseItf interface {
	ProcessPayment(req *dto.NotificationPayment) *res.Err
}

type PaymentUsecase struct {
	env                   *env.Env
	PaymentRepository     repository.PaymentMySQLItf
	TransactionRepository transactionRepo.TransactionMySQLItf
}

func NewPaymentUsecase(env *env.Env, paymentRepository repository.PaymentMySQLItf, transactionRepository transactionRepo.TransactionMySQLItf) PaymentUsecaseItf {
	return &PaymentUsecase{
		env:                   env,
		PaymentRepository:     paymentRepository,
		TransactionRepository: transactionRepository,
	}
}

func (u PaymentUsecase) ProcessPayment(req *dto.NotificationPayment) *res.Err {
	transactionStatusResp := req.TransactionStatus
	fraudStatus := req.FraudStatus

	transactionTime, _ := time.Parse("2006-01-02 15:04:05", req.TransactionTime)
	settlementTime, _ := time.Parse("2006-01-02 15:04:05", req.SettlementTime)

	orderId := strings.SplitN(req.OrderID, "-", 2)

	transactionId, err := uuid.Parse(orderId[0])
	if err != nil {
		return res.ErrBadRequest()
	}

	payment := &entity.Payment{
		TransactionID:     transactionId,
		ReferenceID:       req.ReferenceID,
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
		if err := u.PaymentRepository.Create(payment); err != nil {
			return res.ErrInternalServer()
		}

		transaction.Status = entity.Process
		transaction.Payment = *payment

		if err := u.TransactionRepository.Update(transaction); err != nil {
			return res.ErrInternalServer()
		}
	} else if status == "cancelled" {
		transaction.Status = entity.CancelledBySystem

		if err := u.TransactionRepository.Update(transaction); err != nil {
			return res.ErrInternalServer()
		}
	}

	return nil
}
