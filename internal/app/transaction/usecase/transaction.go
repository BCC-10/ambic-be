package usecase

import (
	notificationRepo "ambic/internal/app/notification/repository"
	PartnerRepo "ambic/internal/app/partner/repository"
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
	GetByUserID(userId uuid.UUID, req dto.GetTransactionByUserIdAndByStatusRequest) (*[]dto.GetTransactionResponse, *res.Err)
	Create(id uuid.UUID, req *dto.CreateTransactionRequest) (string, *res.Err)
	Show(id uuid.UUID) (dto.GetTransactionResponse, *res.Err)
	UpdateStatus(id uuid.UUID, req dto.UpdateTransactionStatusRequest) *res.Err
}

type TransactionUsecase struct {
	db                     *gorm.DB
	env                    *env.Env
	TransactionRepository  repository.TransactionMySQLItf
	ProductRepository      productRepo.ProductMySQLItf
	UserRepository         userRepo.UserMySQLItf
	NotificationRepository notificationRepo.NotificationMySQLItf
	PartnerRepository      PartnerRepo.PartnerMySQLItf
	helper                 helper.HelperIf
	Snap                   midtrans.MidtransIf
}

func NewTransactionUsecase(env *env.Env, db *gorm.DB, transactionRepository repository.TransactionMySQLItf, productRepository productRepo.ProductMySQLItf, userRepository userRepo.UserMySQLItf, notificationRepository notificationRepo.NotificationMySQLItf, partnerRepository PartnerRepo.PartnerMySQLItf, helper helper.HelperIf, snap midtrans.MidtransIf) TransactionUsecaseItf {
	return &TransactionUsecase{
		db:                     db,
		env:                    env,
		TransactionRepository:  transactionRepository,
		ProductRepository:      productRepository,
		UserRepository:         userRepository,
		NotificationRepository: notificationRepository,
		PartnerRepository:      partnerRepository,
		helper:                 helper,
		Snap:                   snap,
	}
}

func (u *TransactionUsecase) GetByUserID(userId uuid.UUID, req dto.GetTransactionByUserIdAndByStatusRequest) (*[]dto.GetTransactionResponse, *res.Err) {
	param := dto.TransactionParam{
		UserID: userId,
	}

	if req.Status != "" {
		param.Status = req.Status
	}

	if req.Limit < 1 {
		req.Limit = u.env.DefaultPaginationLimit
	}

	if req.Page < 1 {
		req.Page = u.env.DefaultPaginationPage
	}

	pagination := dto.PaginationRequest{
		Limit:  req.Limit,
		Page:   req.Page,
		Offset: (req.Page - 1) * req.Limit,
	}

	transactions := new([]entity.Transaction)
	if err := u.TransactionRepository.Get(transactions, param, pagination); err != nil {
		return nil, res.ErrInternalServer()
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

	partnerId, err := uuid.Parse(req.PartnerID)
	if err != nil {
		tx.Rollback()
		return "", res.ErrBadRequest(res.InvalidUUID)
	}

	transaction := &entity.Transaction{
		ID:        transactionId,
		UserID:    userId,
		PartnerID: partnerId,
		Invoice:   u.helper.GenerateInvoiceNumber(),
		Total:     0,
		Status:    entity.WaitingForPayment,
		Note:      req.Note,
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
				return "", res.ErrNotFound(fmt.Sprintf(res.ProductNotFound, item.ProductID))
			}

			tx.Rollback()
			return "", res.ErrInternalServer()
		}

		if product.PartnerID != partnerId {
			tx.Rollback()
			return "", res.ErrBadRequest(fmt.Sprintf(res.ProductNotBelongToPartner, product.Name, product.Partner.Name))
		}

		if product.Stock < uint(item.Qty) {
			tx.Rollback()
			return "", res.ErrBadRequest(fmt.Sprintf(res.InsufficientStock, product.Name))
		}

		transactionDetail := entity.TransactionDetail{
			TransactionID: transactionId,
			ProductID:     productId,
			Qty:           uint(item.Qty),
		}

		partner := new(entity.Partner)
		if err := u.PartnerRepository.Show(partner, dto.PartnerParam{ID: partnerId}); err != nil {
			tx.Rollback()
			return "", res.ErrInternalServer()
		}

		transaction.Total += product.FinalPrice * float32(uint(item.Qty))
		transaction.TransactionDetails = append(transaction.TransactionDetails, transactionDetail)

		items = append(items, dto.TransactionDetail{
			MerchantName: partner.Name,
			ProductID:    item.ProductID,
			Product:      product.ParseDTOGet(nil),
			Qty:          uint(item.Qty),
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

	notification := &entity.Notification{
		UserID:   user.ID,
		Title:    fmt.Sprintf(res.WaitingPaymentTitle),
		Content:  res.WaitingPaymentContent,
		Link:     res.WaitingPaymentLink,
		Button:   res.WaitingPaymentButton,
		PhotoURL: res.WaitingPaymentImageURL,
	}

	if err := u.NotificationRepository.Create(tx, notification); err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer()
	}

	tx.Commit()

	return url, nil
}

func (u *TransactionUsecase) Show(id uuid.UUID) (dto.GetTransactionResponse, *res.Err) {
	transaction := new(entity.Transaction)

	if err := u.TransactionRepository.Show(transaction, dto.TransactionParam{ID: id}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return dto.GetTransactionResponse{}, res.ErrNotFound(res.TransactionNotFound)
		}
	}

	return transaction.ParseDTOGet(), nil
}

func (u *TransactionUsecase) UpdateStatus(id uuid.UUID, req dto.UpdateTransactionStatusRequest) *res.Err {
	tx := u.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transaction := new(entity.Transaction)
	if err := u.TransactionRepository.Show(transaction, dto.TransactionParam{ID: id}); err != nil {
		tx.Rollback()
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.TransactionNotFound)
		}

		return res.ErrInternalServer()
	}

	if transaction.Status == entity.WaitingForPayment || transaction.Status == entity.CancelledBySystem {
		tx.Rollback()
		return res.ErrForbidden(res.NotAllowedToChangeStatus)
	}

	status := new(entity.Status)
	s := entity.Status(req.Status)
	status = &s

	transaction = &entity.Transaction{
		ID:     id,
		Status: *status,
	}

	if err := u.TransactionRepository.Update(tx, transaction); err != nil {
		tx.Rollback()
		return res.ErrInternalServer()
	}

	if req.Status == string(entity.Process) {
		notification := &entity.Notification{
			UserID:   transaction.UserID,
			Title:    res.TransactionProcessTitle,
			Content:  res.TransactionProcessContent,
			Link:     res.TransactionProcessLink,
			Button:   res.TransactionProcessButton,
			PhotoURL: res.TransactionProcessImageURL,
		}

		if err := u.NotificationRepository.Create(tx, notification); err != nil {
			tx.Rollback()
			return res.ErrInternalServer()
		}
	} else if req.Status == string(entity.CancelledBySystem) {
		notification := &entity.Notification{
			UserID:   transaction.UserID,
			Title:    res.TransactionFailedTitle,
			Content:  res.TransactionFailedContent,
			Link:     res.TransactionFailedLink,
			Button:   res.TransactionFailedButton,
			PhotoURL: res.TransactionFailedImageURL,
		}

		if err := u.NotificationRepository.Create(tx, notification); err != nil {
			tx.Rollback()
			return res.ErrInternalServer()
		}
	} else if req.Status == string(entity.Finish) {
		notification := &entity.Notification{
			UserID:   transaction.UserID,
			Title:    res.TransactionFinishTitle,
			Content:  fmt.Sprintf(res.TransactionFinishContent, transaction.Invoice),
			Link:     res.TransactionFinishLink,
			Button:   res.TransactionFinishButton,
			PhotoURL: res.TransactionFinishImageURL,
		}

		if err := u.NotificationRepository.Create(tx, notification); err != nil {
			tx.Rollback()
			return res.ErrInternalServer()
		}
	}

	tx.Commit()

	return nil
}
