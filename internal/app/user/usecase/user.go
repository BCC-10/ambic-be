package usecase

import (
	"ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserUsecaseItf interface {
	UpdateUser(id uuid.UUID, data dto.UpdateUserRequest) *res.Err
}

type UserUsecase struct {
	UserRepository repository.UserMySQLItf
}

func NewUserUsecase(env *env.Env, userRepository repository.UserMySQLItf) UserUsecaseItf {
	return &UserUsecase{
		UserRepository: userRepository,
	}
}

func (u *UserUsecase) UpdateUser(id uuid.UUID, data dto.UpdateUserRequest) *res.Err {
	userDB := new(entity.User)
	if err := u.UserRepository.Get(userDB, dto.UserParam{Id: id}); err != nil {
		return res.ErrNotFound(res.UserNotExists)
	}

	gender := new(entity.Gender)
	if data.Gender != "" {
		g := entity.Gender(data.Gender)
		gender = &g
	}

	user := &entity.User{
		ID:      id,
		Name:    data.Name,
		Phone:   data.Phone,
		Address: data.Address,
		Gender:  gender,
	}

	if data.BornDate != "" {
		bornDate, err := time.Parse("2006-01-02", data.BornDate)
		if err != nil {
			return res.ErrBadRequest(res.InvalidDateFormat)
		}

		user.BornDate = bornDate
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
