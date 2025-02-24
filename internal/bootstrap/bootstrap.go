package bootstrap

import (
	userHandler "ambic/internal/app/user/interface/rest"
	userRepo "ambic/internal/app/user/repository"
	userUsecase "ambic/internal/app/user/usecase"
	"ambic/internal/domain/env"
	"ambic/internal/infra/fiber"
	"ambic/internal/infra/jwt"
	"ambic/internal/infra/mysql"
	"fmt"
	"github.com/go-playground/validator/v10"
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

	err = mysql.Migrate(db)
	if err != nil {
		return err
	}

	val := validator.New()

	jwt := jwt.NewJwt(config)

	app := fiber.New()
	v1 := app.Group("/api/v1")

	userRepository := userRepo.NewUserMySQL(db)
	userUsercase := userUsecase.NewUserUsecase(userRepository, jwt)
	userHandler.NewUserHandler(v1, userUsercase, val)

	return app.Listen(fmt.Sprintf(":%d", config.AppPort))
}
