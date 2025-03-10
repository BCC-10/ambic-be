package bootstrap

import (
	AuthHandler "ambic/internal/app/auth/interface/rest"
	AuthUsecase "ambic/internal/app/auth/usecase"
	BusinessTypeHandler "ambic/internal/app/business_type/interface/rest"
	BusinessTypeRepo "ambic/internal/app/business_type/repository"
	BusinessTypeUsecase "ambic/internal/app/business_type/usecase"
	NotificationHandler "ambic/internal/app/notification/interface/rest"
	NotificationRepo "ambic/internal/app/notification/repository"
	NotificationUsecase "ambic/internal/app/notification/usecase"
	PartnerHandler "ambic/internal/app/partner/interface/rest"
	PartnerRepo "ambic/internal/app/partner/repository"
	PartnerUsecase "ambic/internal/app/partner/usecase"
	PaymentHandler "ambic/internal/app/payment/interface/rest"
	PaymentRepo "ambic/internal/app/payment/repository"
	PaymentUsecase "ambic/internal/app/payment/usecase"
	ProductHandler "ambic/internal/app/product/interface/rest"
	ProductRepo "ambic/internal/app/product/repository"
	ProductUsecase "ambic/internal/app/product/usecase"
	RatingHandler "ambic/internal/app/rating/interface/rest"
	RatingRepo "ambic/internal/app/rating/repository"
	RatingUsecase "ambic/internal/app/rating/usecase"
	TransactionHandler "ambic/internal/app/transaction/interface/rest"
	TransactionRepo "ambic/internal/app/transaction/repository"
	TransactionUsecase "ambic/internal/app/transaction/usecase"
	UserHandler "ambic/internal/app/user/interface/rest"
	UserRepo "ambic/internal/app/user/repository"
	UserUsecase "ambic/internal/app/user/usecase"
	"ambic/internal/domain/env"
	"ambic/internal/infra/code"
	"ambic/internal/infra/email"
	"ambic/internal/infra/fiber"
	"ambic/internal/infra/helper"
	"ambic/internal/infra/jwt"
	"ambic/internal/infra/limiter"
	"ambic/internal/infra/maps"
	"ambic/internal/infra/midtrans"
	"ambic/internal/infra/mysql"
	"ambic/internal/infra/oauth"
	"ambic/internal/infra/redis"
	"ambic/internal/infra/supabase"
	"ambic/internal/middleware"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Start() error {
	config, err := env.New()
	if err != nil {
		panic(err)
	}

	db, err := mysql.New(fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	))

	if err := mysql.Migrate(db); err != nil {
		return err
	}

	h := helper.New(config)

	v := validator.New()

	j := jwt.NewJwt(config)

	m := middleware.NewMiddleware(j)

	r := redis.NewRedis(config)

	c := code.NewCode(config)

	e := email.NewEmail(config)

	o := oauth.NewOAuth(config)

	s := supabase.New(config)

	ma := maps.NewMaps(config)

	snap := midtrans.New(config)

	app := fiber.New(config)
	app.Get("/metrics", monitor.New())
	v1 := app.Group("/api/v1")

	l := limiter.NewLimiter(r)

	businessTypeRepository := BusinessTypeRepo.NewBusinessTypeMySQL(db)
	paymentRepository := PaymentRepo.NewPaymentMySQL(db)
	productRepository := ProductRepo.NewProductMySQL(db)
	partnerRepository := PartnerRepo.NewPartnerMySQL(db)
	ratingRepository := RatingRepo.NewRatingMySQL(db)
	transactionRepository := TransactionRepo.NewTransactionMySQL(db)
	userRepository := UserRepo.NewUserMySQL(db)
	notificationRepository := NotificationRepo.NewNotificationMySQL(db)

	notificationUsecase := NotificationUsecase.NewNotificationUsecase(config, notificationRepository)
	NotificationHandler.NewNotificationHandler(v1, notificationUsecase, m)

	userUsecase := UserUsecase.NewUserUsecase(config, userRepository, s, h)
	UserHandler.NewUserHandler(v1, userUsecase, v, m, h)

	authUsecase := AuthUsecase.NewAuthUsecase(config, db, userRepository, notificationRepository, j, c, e, r, o)
	AuthHandler.NewAuthHandler(v1, authUsecase, v, l)

	businessTypeUsecase := BusinessTypeUsecase.NewBusinessTypeUsecase(config, businessTypeRepository)
	BusinessTypeHandler.NewBusinessTypeHandler(v1, businessTypeUsecase, m)

	productUsecase := ProductUsecase.NewProductUsecase(config, db, productRepository, s, h)
	ProductHandler.NewProductHandler(v1, productUsecase, v, m, h)

	partnerUsecase := PartnerUsecase.NewPartnerUsecase(config, partnerRepository, userRepository, businessTypeRepository, productRepository, ratingRepository, transactionRepository, s, h, ma, j)
	PartnerHandler.NewPartnerHandler(v1, partnerUsecase, v, m, h)

	ratingUsecase := RatingUsecase.NewRatingUsecase(config, ratingRepository, productRepository, transactionRepository, s, h)
	RatingHandler.NewRatingHandler(v1, ratingUsecase, v, m, h)

	transactionUsecase := TransactionUsecase.NewTransactionUsecase(config, db, transactionRepository, productRepository, userRepository, notificationRepository, partnerRepository, h, snap)
	TransactionHandler.NewTransactionHandler(v1, transactionUsecase, v, m)

	paymentUsecase := PaymentUsecase.NewPaymentUsecase(config, paymentRepository, transactionRepository)
	PaymentHandler.NewPaymentHandler(v1, paymentUsecase, v)

	return app.Listen(fmt.Sprintf(":%d", config.AppPort))
}
