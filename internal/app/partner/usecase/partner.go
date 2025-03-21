package usecase

import (
	businessTypeRepo "ambic/internal/app/business_type/repository"
	"ambic/internal/app/partner/repository"
	productRepo "ambic/internal/app/product/repository"
	ratingRepo "ambic/internal/app/rating/repository"
	transactionRepo "ambic/internal/app/transaction/repository"
	userRepo "ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/code"
	"ambic/internal/infra/email"
	"ambic/internal/infra/helper"
	"ambic/internal/infra/jwt"
	"ambic/internal/infra/maps"
	"ambic/internal/infra/mysql"
	"ambic/internal/infra/redis"
	res "ambic/internal/infra/response"
	"ambic/internal/infra/supabase"
	"ambic/internal/infra/telegram"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
)

type PartnerUsecaseItf interface {
	ShowPartner(id uuid.UUID) (dto.GetPartnerResponse, *res.Err)
	RegisterPartner(id uuid.UUID, data dto.RegisterPartnerRequest) (string, *res.Err)
	VerifyPartner(request dto.VerifyPartnerRequest) (string, *res.Err)
	GetProducts(id uuid.UUID, pagination dto.PaginationRequest) ([]dto.GetProductResponse, *dto.PaginationResponse, *res.Err)
	UpdatePhoto(id uuid.UUID, data dto.UpdatePhotoRequest) *res.Err
	GetStatistics(id uuid.UUID) (dto.GetPartnerStatisticResponse, *res.Err)
	GetTransactions(id uuid.UUID, data dto.GetPartnerTransactionRequest) ([]dto.GetTransactionResponse, *dto.PaginationResponse, *res.Err)
	RequestPartnerVerification(data dto.RequestPartnerVerificationRequest) *res.Err
}

type PartnerUsecase struct {
	env                    *env.Env
	PartnerRepository      repository.PartnerMySQLItf
	UserRepository         userRepo.UserMySQLItf
	BusinessTypeRepository businessTypeRepo.BusinessTypeMySQLItf
	ProductRepository      productRepo.ProductMySQLItf
	TransactionRepository  transactionRepo.TransactionMySQLItf
	RatingRepository       ratingRepo.RatingMySQLItf
	Maps                   maps.MapsIf
	Supabase               supabase.SupabaseIf
	code                   code.CodeIf
	helper                 helper.HelperIf
	jwt                    jwt.JWTIf
	redis                  redis.RedisIf
	email                  email.EmailIf
	telegram               telegram.TelegramIf
	db                     *gorm.DB
}

func NewPartnerUsecase(env *env.Env, db *gorm.DB, partnerRepository repository.PartnerMySQLItf, userRepository userRepo.UserMySQLItf, businessTypeRepository businessTypeRepo.BusinessTypeMySQLItf, productRepository productRepo.ProductMySQLItf, ratingRepository ratingRepo.RatingMySQLItf, transactionRepository transactionRepo.TransactionMySQLItf, supabase supabase.SupabaseIf, helper helper.HelperIf, maps maps.MapsIf, jwt jwt.JWTIf, email email.EmailIf, code code.CodeIf, redis redis.RedisIf, telegram telegram.TelegramIf) PartnerUsecaseItf {
	return &PartnerUsecase{
		env:                    env,
		PartnerRepository:      partnerRepository,
		ProductRepository:      productRepository,
		Maps:                   maps,
		RatingRepository:       ratingRepository,
		UserRepository:         userRepository,
		BusinessTypeRepository: businessTypeRepository,
		TransactionRepository:  transactionRepository,
		Supabase:               supabase,
		helper:                 helper,
		jwt:                    jwt,
		code:                   code,
		redis:                  redis,
		email:                  email,
		telegram:               telegram,
		db:                     db,
	}
}

func (u *PartnerUsecase) RegisterPartner(id uuid.UUID, data dto.RegisterPartnerRequest) (string, *res.Err) {
	tx := u.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := new(entity.User)
	if err := u.UserRepository.Show(user, dto.UserParam{Id: id}); err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer()
	}

	if user.Name == "" || user.Phone == "" || user.Address == "" || user.Gender == nil {
		tx.Rollback()
		return "", res.ErrForbidden(res.ProfileNotFilledCompletely)
	}

	if data.Instagram[0] == '@' {
		data.Instagram = data.Instagram[1:]
	}

	businessTypeId, err := uuid.Parse(data.BusinessTypeID)
	if err != nil {
		tx.Rollback()
		return "", res.ErrBadRequest(res.InvalidBusinessType)
	}

	businessType := new(entity.BusinessType)
	if err := u.BusinessTypeRepository.Show(businessType, dto.BusinessTypeParam{ID: businessTypeId}); err != nil {
		tx.Rollback()
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return "", res.ErrBadRequest(res.InvalidBusinessType)
		}

		return "", res.ErrInternalServer()
	}

	placeDetails, err := u.Maps.GetPlaceDetails(data.PlaceID)
	if err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer(err.Error())
	}

	partner := entity.Partner{
		UserID:         id,
		Name:           data.Name,
		Address:        data.Address,
		City:           data.City,
		Instagram:      data.Instagram,
		PlaceID:        data.PlaceID,
		Latitude:       placeDetails.Lat,
		Longitude:      placeDetails.Long,
		BusinessTypeID: businessTypeId,
		PhotoURL:       u.env.DefaultPartnerProfilePhotoURL,
	}

	if data.Photo != nil {
		if err := u.helper.ValidateImage(data.Photo); err != nil {
			tx.Rollback()
			return "", err
		}

		src, err := data.Photo.Open()
		if err != nil {
			tx.Rollback()
			return "", res.ErrInternalServer()
		}

		defer src.Close()

		bucket := u.env.SupabaseBucket
		path := "partners/" + uuid.NewString() + filepath.Ext(data.Photo.Filename)
		contentType := data.Photo.Header.Get("Content-Type")

		photoURL, err := u.Supabase.UploadFile(bucket, path, contentType, src)
		if err != nil {
			return "", res.ErrInternalServer()
		}

		partner.PhotoURL = photoURL
	}

	if err := u.PartnerRepository.Create(tx, &partner); err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer()
	}

	token, err := u.jwt.GenerateToken(user.ID, user.IsVerified, user.Partner.ID, user.Partner.IsVerified)
	if err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer()
	}

	if err := u.email.SendPartnerRegistrationEmail(user.Email, partner.Name); err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer()
	}

	msg := dto.PartnerRegistrationTelegramMessage{
		UserID:               user.ID.String(),
		UserName:             user.Name,
		UserUsername:         user.Username,
		UserEmail:            user.Email,
		UserPhone:            user.Phone,
		UserAddress:          user.Address,
		UserGender:           user.Gender.String(),
		UserRegisteredAt:     user.CreatedAt.String(),
		PartnerID:            partner.ID.String(),
		BusinessType:         businessType.Name,
		BusinessName:         partner.Name,
		BusinessAddress:      partner.Address,
		BusinessCity:         partner.City,
		BusinessGmaps:        u.Maps.GenerateGoogleMapsURL(partner.PlaceID),
		BusinessInstagram:    partner.Instagram,
		BusinessPhoto:        partner.PhotoURL,
		BusinessRegisteredAt: partner.CreatedAt.String(),
	}

	if err := u.telegram.SendPartnerRegistrationMessage(msg); err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer(err.Error())
	}

	tx.Commit()

	return token, nil
}

func (u *PartnerUsecase) RequestPartnerVerification(data dto.RequestPartnerVerificationRequest) *res.Err {
	if data.Token != u.env.PartnerVerificationToken {
		return res.ErrForbidden(res.InvalidVerificationToken)
	}

	user := new(entity.User)
	if err := u.UserRepository.Show(user, dto.UserParam{Email: data.Email}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound(res.PartnerNotExists)
		}

		return res.ErrInternalServer()
	}

	if user.Partner.ID == uuid.Nil {
		return res.ErrNotFound(res.PartnerNotExists)
	}

	if user.Partner.IsVerified {
		return res.ErrForbidden(res.PartnerVerified)
	}

	token, err := u.code.GenerateToken()
	if err != nil {
		return res.ErrInternalServer()
	}

	if err := u.redis.Set("p"+data.Email, []byte(token), u.env.PartnerVerificationTokenExpires); err != nil {
		return res.ErrInternalServer()
	}

	if err := u.email.SendPartnerVerificationEmail(user.Email, user.Partner.Name, token); err != nil {
		return res.ErrInternalServer()
	}

	return nil
}

func (u *PartnerUsecase) VerifyPartner(data dto.VerifyPartnerRequest) (string, *res.Err) {
	tx := u.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := new(entity.User)
	if err := u.UserRepository.Show(user, dto.UserParam{Email: data.Email}); err != nil {
		tx.Rollback()
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return "", res.ErrNotFound(res.PartnerNotExists)
		}

		return "", res.ErrInternalServer()
	}

	if user.Partner.ID == uuid.Nil {
		tx.Rollback()
		return "", res.ErrNotFound(res.PartnerNotExists)
	}

	if user.Partner.IsVerified {
		tx.Rollback()
		return "", res.ErrForbidden(res.PartnerVerified)
	}

	token, err := u.redis.Get("p" + data.Email)
	if err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer()
	}

	if string(token) != data.Token {
		tx.Rollback()
		return "", res.ErrForbidden(res.InvalidVerificationToken)
	}

	user.Partner.IsVerified = true

	if err := u.PartnerRepository.Update(tx, &user.Partner); err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer()
	}

	newJWTToken, err := u.jwt.GenerateToken(user.ID, user.IsVerified, user.Partner.ID, user.Partner.IsVerified)
	if err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer()
	}

	if err := u.redis.Delete("p" + data.Email); err != nil {
		tx.Rollback()
		return "", res.ErrInternalServer()
	}

	tx.Commit()

	return newJWTToken, nil
}

func (u *PartnerUsecase) GetProducts(id uuid.UUID, pagination dto.PaginationRequest) ([]dto.GetProductResponse, *dto.PaginationResponse, *res.Err) {
	pagination = u.helper.CreatePagination(pagination)

	products := new([]entity.Product)
	if err := u.ProductRepository.GetByPartnerId(products, dto.ProductParam{PartnerId: id}, pagination); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return nil, nil, res.ErrNotFound(res.PartnerNotExists)
		}

		return nil, nil, res.ErrInternalServer()
	}

	productsResponse := make([]dto.GetProductResponse, 0)
	for _, product := range *products {
		productsResponse = append(productsResponse, product.ParseDTOGet(nil))
	}

	totalProducts, err := u.ProductRepository.GetTotalProductsByPartnerId(id)
	if err != nil {
		return nil, nil, res.ErrInternalServer()
	}

	pg := u.helper.CalculatePagination(pagination, totalProducts)

	return productsResponse, &pg, nil
}

func (u *PartnerUsecase) ShowPartner(id uuid.UUID) (dto.GetPartnerResponse, *res.Err) {
	partner := new(entity.Partner)
	if err := u.PartnerRepository.Show(partner, dto.PartnerParam{ID: id}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return dto.GetPartnerResponse{}, res.ErrNotFound(res.PartnerNotExists)
		}

		return dto.GetPartnerResponse{}, res.ErrInternalServer()
	}

	return partner.ParseDTOGet(), nil
}

func (u *PartnerUsecase) UpdatePhoto(id uuid.UUID, data dto.UpdatePhotoRequest) *res.Err {
	tx := u.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	partnerDB := new(entity.Partner)
	if err := u.PartnerRepository.Show(partnerDB, dto.PartnerParam{ID: id}); err != nil {
		return res.ErrInternalServer()
	}

	if err := u.helper.ValidateImage(data.Photo); err != nil {
		return err
	}

	src, err := data.Photo.Open()
	if err != nil {
		return res.ErrInternalServer()
	}

	defer src.Close()

	bucket := u.env.SupabaseBucket
	path := "partners/" + uuid.NewString() + filepath.Ext(data.Photo.Filename)
	contentType := data.Photo.Header.Get("Content-Type")

	photoURL, err := u.Supabase.UploadFile(bucket, path, contentType, src)
	if err != nil {
		return res.ErrInternalServer()
	}

	partner := &entity.Partner{
		ID:       id,
		PhotoURL: photoURL,
	}

	if err := u.PartnerRepository.Update(tx, partner); err != nil {
		return res.ErrInternalServer()
	}

	if partnerDB.PhotoURL != u.env.DefaultPartnerProfilePhotoURL {
		oldPhotoURL := partnerDB.PhotoURL
		index := strings.Index(oldPhotoURL, bucket)
		oldPhotoPath := oldPhotoURL[index+len(bucket+"/"):]

		if err = u.Supabase.DeleteFile(bucket, oldPhotoPath); err != nil {
			tx.Rollback()
			return res.ErrInternalServer()
		}
	}

	tx.Commit()

	return nil
}

func (u *PartnerUsecase) GetStatistics(id uuid.UUID) (dto.GetPartnerStatisticResponse, *res.Err) {
	resp := new(dto.GetPartnerStatisticResponse)

	partner := new(entity.Partner)
	if err := u.PartnerRepository.Show(partner, dto.PartnerParam{ID: id}); err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return *resp, res.ErrNotFound(res.PartnerNotExists)
		}

		return *resp, res.ErrInternalServer()
	}

	totalRatings, err := u.RatingRepository.GetTotalRatingsByPartnerId(id)
	if err != nil {
		return *resp, res.ErrInternalServer()
	}

	totalProducts, err := u.ProductRepository.GetTotalProductsByPartnerId(id)
	if err != nil {
		return *resp, res.ErrInternalServer()
	}

	if err := u.PartnerRepository.ShowWithTransactions(partner, dto.PartnerParam{ID: id}); err != nil {
		return *resp, res.ErrInternalServer()
	}

	var totalTransactions int64
	var totalRevenue float32
	for _, transaction := range partner.Transactions {
		if transaction.Status == entity.Finish {
			totalTransactions++
			totalRevenue += transaction.Total
		}
	}

	resp.TotalRatings = totalRatings
	resp.TotalProducts = totalProducts
	resp.TotalTransactions = totalTransactions
	resp.TotalRevenue = totalRevenue

	return *resp, nil
}

func (u *PartnerUsecase) GetTransactions(id uuid.UUID, data dto.GetPartnerTransactionRequest) ([]dto.GetTransactionResponse, *dto.PaginationResponse, *res.Err) {
	pagination := u.helper.CreatePagination(dto.PaginationRequest{
		Limit: data.Limit,
		Page:  data.Page,
	})

	param := dto.TransactionParam{PartnerID: id}

	if data.Status != "" {
		param.Status = data.Status
	}

	transactions := new([]entity.Transaction)
	totalTransactions, err := u.TransactionRepository.Get(transactions, param, pagination)
	if err != nil {
		if mysql.CheckError(err, gorm.ErrRecordNotFound) {
			return nil, nil, res.ErrNotFound(res.PartnerNotExists)
		}

		return nil, nil, res.ErrInternalServer()
	}

	transactionsResponse := make([]dto.GetTransactionResponse, 0)
	for _, transaction := range *transactions {
		if transaction.Status == entity.CancelledBySystem {
			continue
		}

		transactionsResponse = append(transactionsResponse, transaction.ParseDTOGet())
	}

	pg := u.helper.CalculatePagination(pagination, totalTransactions)

	return transactionsResponse, &pg, nil
}
