package usecase

import (
	"ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/code"
	"ambic/internal/infra/email"
	"ambic/internal/infra/jwt"
	"ambic/internal/infra/redis"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseItf interface {
	Register(dto.Register) error
	Login(login dto.Login) (string, error)
	RequestOTP(requestOTP dto.RequestOTP) error
	VerifyOTP(verifyOTP dto.VerifyOTP) error
}

type UserUsecase struct {
	UserRepository repository.UserMySQLItf
	jwt            jwt.JWTIf
	code           code.CodeIf
	email          email.EmailIf
	redis          redis.RedisIf
	env            *env.Env
}

func NewUserUsecase(env *env.Env, userRepository repository.UserMySQLItf, jwt jwt.JWTIf, code code.CodeIf, email email.EmailIf, redis redis.RedisIf) UserUsecaseItf {
	return &UserUsecase{
		UserRepository: userRepository,
		jwt:            jwt,
		code:           code,
		email:          email,
		redis:          redis,
		env:            env,
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

func (u *UserUsecase) RequestOTP(data dto.RequestOTP) error {
	user := new(entity.User)
	err := u.UserRepository.Get(user, dto.UserParam{Email: data.Email})
	if err != nil {
		return errors.New("email is not registered")
	}

	otp, err := u.code.GenerateOTP()
	if err != nil {
		return errors.New("failed to generate OTP")
	}

	err = u.redis.Set(data.Email, []byte(otp), u.env.OTPExpiresTime)
	if err != nil {
		return errors.New("failed to save OTP")
	}

	err = u.email.SendOTP(data.Email, otp)
	if err != nil {
		return errors.New("failed to send OTP")
	}

	return nil
}

func (u *UserUsecase) VerifyOTP(data dto.VerifyOTP) error {
	savedOTP, err := u.redis.Get(data.Email)
	if err != nil {
		return errors.New("OTP is expired or invalid")
	}

	if string(savedOTP) != data.OTP {
		return errors.New("OTP is expired or invalid")
	}

	err = u.redis.Set(data.Email, nil, 0)
	if err != nil {
		return errors.New("failed to delete OTP")
	}

	err = u.UserRepository.Activate(&entity.User{Email: data.Email})
	if err != nil {
		return errors.New("failed to activate user")
	}

	return nil
}
