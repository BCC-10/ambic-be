package usecase

import (
	"ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseItf interface {
	UpdateUser(id uuid.UUID, data dto.UpdateUser) *res.Err
}

type UserUsecase struct {
	UserRepository repository.UserMySQLItf
}

func NewUserUsecase(env *env.Env, userRepository repository.UserMySQLItf) UserUsecaseItf {
	return &UserUsecase{
		UserRepository: userRepository,
	}
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
