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
	res "ambic/internal/infra/response"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseItf interface {
	Register(dto.Register) *res.Err
	Login(login dto.Login) (string, *res.Err)
	RequestOTP(requestOTP dto.RequestOTP) *res.Err
	VerifyUser(verifyUser dto.VerifyOTP) *res.Err
	ForgotPassword(resetPassword dto.ForgotPassword) *res.Err
	ResetPassword(data dto.ResetPassword) *res.Err
	UpdateUser(id uuid.UUID, data dto.UpdateUser) *res.Err
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

func (u *UserUsecase) Register(data dto.Register) *res.Err {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return res.ErrInternalServer()
	}

	user := entity.User{
		ID:         uuid.New(),
		Name:       data.Name,
		Username:   data.Username,
		Email:      data.Email,
		Password:   string(hashedPassword),
		IsVerified: false,
	}

	var dbUser entity.User
	if err := u.UserRepository.Get(&dbUser, dto.UserParam{Email: user.Email}); err == nil {
		return res.ErrBadRequest(res.EmailExist)
	}

	if err := u.UserRepository.Get(&dbUser, dto.UserParam{Username: user.Username}); err == nil {
		return res.ErrBadRequest(res.UsernameExist)
	}

	if err := u.UserRepository.Create(&user); err != nil {
		return res.ErrInternalServer()
	}

	return nil
}

func (u *UserUsecase) Login(data dto.Login) (string, *res.Err) {
	user := new(entity.User)

	err := u.UserRepository.Check(user, data)
	if err != nil {
		return "", res.ErrUnauthorized(res.IncorrectIdentifier)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return "", res.ErrUnauthorized(res.IncorrectIdentifier)
	}

	if !user.IsVerified {
		return "", res.ErrForbidden(res.UserNotVerified)
	}

	token, err := u.jwt.GenerateToken(user.ID, user.IsVerified)
	if err != nil {
		return "", res.ErrInternalServer()
	}

	return token, nil
}

func (u *UserUsecase) RequestOTP(data dto.RequestOTP) *res.Err {
	user := new(entity.User)
	err := u.UserRepository.Get(user, dto.UserParam{Email: data.Email})
	if err != nil {
		return res.ErrNotFound(res.UserNotExists)
	}

	otp, err := u.code.GenerateOTP()
	if err != nil {
		return res.ErrInternalServer()
	}

	err = u.redis.Set(data.Email, []byte(otp), u.env.OTPExpiresTime)
	if err != nil {
		return res.ErrInternalServer()
	}

	err = u.email.SendOTP(data.Email, otp)
	if err != nil {
		return res.ErrInternalServer()
	}

	return nil
}

func (u *UserUsecase) VerifyUser(data dto.VerifyOTP) *res.Err {
	user := new(entity.User)
	err := u.UserRepository.Get(user, dto.UserParam{Email: data.Email})
	if err != nil {
		return res.ErrNotFound(res.UserNotExists)
	}

	if user.IsVerified {
		return res.ErrForbidden(res.UserVerified)
	}

	savedOTP, err := u.redis.Get(data.Email)
	if err != nil {
		return res.ErrInternalServer()
	}

	if string(savedOTP) != data.OTP {
		return res.ErrBadRequest(res.InvalidOTP)
	}

	err = u.UserRepository.Verify(user)
	if err != nil {
		return res.ErrInternalServer()
	}

	err = u.redis.Delete(data.Email)
	if err != nil {
		return res.ErrInternalServer()
	}

	return nil
}

func (u *UserUsecase) ForgotPassword(data dto.ForgotPassword) *res.Err {
	user := new(entity.User)
	if err := u.UserRepository.Get(user, dto.UserParam{Email: data.Email}); err != nil {
		return res.ErrNotFound(res.UserNotExists)
	}

	if !user.IsVerified {
		return res.ErrForbidden(res.UserNotVerified)
	}

	token, err := u.code.GenerateToken()
	if err != nil {
		return res.ErrInternalServer()
	}

	err = u.redis.Set(user.Email, []byte(token), u.env.TokenExpiresTime)
	if err != nil {
		return res.ErrInternalServer()
	}

	err = u.email.SendResetPasswordLink(user.Email, token)
	if err != nil {
		return res.ErrInternalServer()
	}

	return nil
}

func (u *UserUsecase) ResetPassword(data dto.ResetPassword) *res.Err {
	user := new(entity.User)
	if err := u.UserRepository.Get(user, dto.UserParam{Email: data.Email}); err != nil {
		return res.ErrNotFound(res.UserNotExists)
	}

	if !user.IsVerified {
		return res.ErrForbidden(res.UserNotVerified)
	}

	savedToken, err := u.redis.Get(data.Email)
	if err != nil {
		return res.ErrInternalServer()
	}

	if string(savedToken) != data.Token {
		return res.ErrBadRequest(res.InvalidToken)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return res.ErrInternalServer()
	}

	user.Password = string(hashedPassword)
	err = u.UserRepository.Update(user)
	if err != nil {
		return res.ErrInternalServer()
	}

	err = u.redis.Delete(data.Email)
	if err != nil {
		return res.ErrInternalServer()
	}

	return nil
}

func (u *UserUsecase) UpdateUser(id uuid.UUID, data dto.UpdateUser) *res.Err {
	userDB := new(entity.User)
	if err := u.UserRepository.Get(userDB, dto.UserParam{Id: id}); err != nil {
		return res.ErrNotFound(res.UserNotExists)
	}

	user := &entity.User{
		ID:   id,
		Name: data.Name,
	}

	if data.NewPassword != "" {
		err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(data.OldPassword))
		if err != nil {
			return res.ErrForbidden(res.IncorrectOldPassword)
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return res.ErrInternalServer()
		}

		user.Password = string(hashedPassword)
	}

	err := u.UserRepository.Update(user)
	if err != nil {
		return res.ErrInternalServer()
	}

	return nil
}
