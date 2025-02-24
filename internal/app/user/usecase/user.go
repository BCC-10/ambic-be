package usecase

import (
	"ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseItf interface {
	Register(dto.Register) error
}

type UserUsecase struct {
	UserRepository repository.UserMySQLItf
}

func NewUserUsecase(userRepository repository.UserMySQLItf) UserUsecaseItf {
	return &UserUsecase{
		UserRepository: userRepository,
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
