package bootstrap

import (
	userHandler "ambic/internal/app/user/interface/rest"
	userRepo "ambic/internal/app/user/repository"
	userUsecase "ambic/internal/app/user/usecase"
	"ambic/internal/domain/env"
	"ambic/internal/infra/code"
	"ambic/internal/infra/email"
	"ambic/internal/infra/fiber"
	"ambic/internal/infra/jwt"
	"ambic/internal/infra/limiter"
	"ambic/internal/infra/mysql"
	"ambic/internal/infra/redis"
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

	err = mysql.Migrate(db)
	if err != nil {
		return err
	}

	v := validator.New()

	j := jwt.NewJwt(config)

	m := middleware.NewMiddleware(j)

	r := redis.NewRedis(config)

	c := code.NewCode(config)

	e := email.NewEmail(config)

	app := fiber.New()
	app.Get("/metrics", monitor.New())
	v1 := app.Group("/api/v1")

	l := limiter.NewLimiter(r)

	userRepository := userRepo.NewUserMySQL(db)
	userUsercase := userUsecase.NewUserUsecase(config, userRepository, j, c, e, r)
	userHandler.NewUserHandler(v1, userUsercase, v, m, l)

	return app.Listen(fmt.Sprintf(":%d", config.AppPort))
}
