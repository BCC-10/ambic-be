package usecase

import (
	"ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	"ambic/internal/infra/code"
	"ambic/internal/infra/email"
	"ambic/internal/infra/jwt"
	"ambic/internal/infra/oauth"
	"ambic/internal/infra/redis"
	res "ambic/internal/infra/response"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseItf interface {
	Register(dto.RegisterRequest) *res.Err
	Login(login dto.LoginRequest) (string, *res.Err)
	RequestOTP(requestOTP dto.OTPRequest) *res.Err
	VerifyUser(verifyUser dto.VerifyOTPRequest) *res.Err
	ForgotPassword(resetPassword dto.ForgotPasswordRequest) *res.Err
	ResetPassword(data dto.ResetPasswordRequest) *res.Err
	GoogleLogin() (string, *res.Err)
	GoogleCallback(code string, state string) (string, *res.Err)
}

type AuthUsecase struct {
	UserRepository repository.UserMySQLItf
	jwt            jwt.JWTIf
	code           code.CodeIf
	email          email.EmailIf
	redis          redis.RedisIf
	env            *env.Env
	OAuth          oauth.OAuthIf
}

func NewAuthUsecase(env *env.Env, userRepository repository.UserMySQLItf, jwt jwt.JWTIf, code code.CodeIf, email email.EmailIf, redis redis.RedisIf, oauth oauth.OAuthIf) AuthUsecaseItf {
	return &AuthUsecase{
		UserRepository: userRepository,
		jwt:            jwt,
		code:           code,
		email:          email,
		redis:          redis,
		env:            env,
		OAuth:          oauth,
	}
}

func (u *AuthUsecase) Register(data dto.RegisterRequest) *res.Err {
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

func (u *AuthUsecase) Login(data dto.LoginRequest) (string, *res.Err) {
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

func (u *AuthUsecase) RequestOTP(data dto.OTPRequest) *res.Err {
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

func (u *AuthUsecase) VerifyUser(data dto.VerifyOTPRequest) *res.Err {
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

func (u *AuthUsecase) ForgotPassword(data dto.ForgotPasswordRequest) *res.Err {
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

func (u *AuthUsecase) ResetPassword(data dto.ResetPasswordRequest) *res.Err {
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

func (u *AuthUsecase) GoogleLogin() (string, *res.Err) {
	state, err := u.code.GenerateToken()
	if err != nil {
		return "", res.ErrInternalServer()
	}

	err = u.redis.Set(state, []byte(state), u.env.StateExpiresTime)
	if err != nil {
		return "", res.ErrInternalServer()
	}

	url, err := u.OAuth.GenerateGoogleAuthLink(state)
	if err != nil {
		return "", res.ErrInternalServer()
	}

	return url, nil
}

func (u *AuthUsecase) GoogleCallback(code string, state string) (string, *res.Err) {
	savedState, err := u.redis.Get(state)
	if err != nil {
		return "", res.ErrInternalServer()
	}

	if string(savedState) != state {
		return "", res.ErrBadRequest(res.InvalidState)
	}

	token, err := u.OAuth.ExchangeToken(code)
	if err != nil {
		return "", res.ErrInternalServer()
	}

	profile, err := u.OAuth.GetUserProfile(token)
	if err != nil {
		return "", res.ErrInternalServer()
	}

	user := &entity.User{
		ID:         uuid.New(),
		Name:       profile.Name,
		Username:   profile.Username,
		Email:      profile.Email,
		IsVerified: profile.IsVerified,
	}

	var dbUser entity.User
	err = u.UserRepository.Get(&dbUser, dto.UserParam{Email: user.Email})
	if err != nil {
		err = u.UserRepository.Create(user)
		if err != nil {
			return "", res.ErrInternalServer()
		}
	} else {
		user = &dbUser
	}

	jwtToken, err := u.jwt.GenerateToken(user.ID, user.IsVerified)
	if err != nil {
		return "", res.ErrInternalServer()
	}

	err = u.redis.Delete(state)
	if err != nil {
		return "", res.ErrInternalServer()
	}

	return jwtToken, nil
}
