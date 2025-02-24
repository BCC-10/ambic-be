package usecase

import (
	"ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/infra/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseItf interface {
	Register(dto.Register) error
	Login(login dto.Login) (string, error)
}

type UserUsecase struct {
	UserRepository repository.UserMySQLItf
	jwt            jwt.JWTIf
}

func NewUserUsecase(userRepository repository.UserMySQLItf, jwt jwt.JWTIf) UserUsecaseItf {
	return &UserUsecase{
		UserRepository: userRepository,
		jwt:            jwt,
	}
}

func (u *UserUsecase) Register(register dto.Register) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), 10)
	if err != nil {
		return err
	}

	user := entity.User{
		ID:       uuid.New(),
		Name:     register.Name,
		Username: register.Username,
		Email:    register.Email,
		Password: string(hashedPassword),
	}

	return u.UserRepository.Create(&user)
}

func (u *UserUsecase) Login(login dto.Login) (string, error) {
	user := new(entity.User)

	err := u.UserRepository.Get(user, dto.UserParam{Email: login.Email})
	if err != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "email or password is incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return "", err
	}

	token, err := u.jwt.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
