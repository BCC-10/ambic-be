package usecase

import (
	"ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/infra/code"
	"ambic/internal/infra/email"
	"ambic/internal/infra/jwt"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseItf interface {
	Register(dto.Register) error
	Login(login dto.Login) (string, error)
	RequestOTP(requestOTP dto.RequestOTP) error
}

type UserUsecase struct {
	UserRepository repository.UserMySQLItf
	jwt            jwt.JWTIf
	code           code.CodeIf
	email          email.EmailIf
}

func NewUserUsecase(userRepository repository.UserMySQLItf, jwt jwt.JWTIf, code code.CodeIf, email email.EmailIf) UserUsecaseItf {
	return &UserUsecase{
		UserRepository: userRepository,
		jwt:            jwt,
		code:           code,
		email:          email,
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
		IsActive: false,
	}

	return u.UserRepository.Create(&user)
}

func (u *UserUsecase) Login(login dto.Login) (string, error) {
	user := new(entity.User)

	err := u.UserRepository.Get(user, dto.UserParam{Email: login.Email})
	if err != nil {
		return "", errors.New("email or password is incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return "", errors.New("email or password is incorrect")
	}

	if !user.IsActive {
		return "", errors.New("user is not active, check your email to activate your account")
	}

	token, err := u.jwt.GenerateToken(user.ID, user.IsActive)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUsecase) RequestOTP(requestOTP dto.RequestOTP) error {
	otp, err := u.code.GenerateOTP(6)
	if err != nil {
		return err
	}

	return u.email.SendOTP(requestOTP.Email, otp)
}
