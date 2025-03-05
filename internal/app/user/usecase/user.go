package usecase

import (
	"ambic/internal/app/user/repository"
	"ambic/internal/domain/dto"
	"ambic/internal/domain/entity"
	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
	"ambic/internal/infra/supabase"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"path/filepath"
	"strings"
	"time"
)

type UserUsecaseItf interface {
	UpdateUser(id uuid.UUID, data dto.UpdateUserRequest) *res.Err
}

type UserUsecase struct {
	UserRepository repository.UserMySQLItf
	Supabase       supabase.SupabaseIf
	env            *env.Env
}

func NewUserUsecase(env *env.Env, userRepository repository.UserMySQLItf, supabase supabase.SupabaseIf) UserUsecaseItf {
	return &UserUsecase{
		UserRepository: userRepository,
		Supabase:       supabase,
		env:            env,
	}
}

func (u *UserUsecase) UpdateUser(id uuid.UUID, data dto.UpdateUserRequest) *res.Err {
	userDB := new(entity.User)
	if err := u.UserRepository.Show(userDB, dto.UserParam{Id: id}); err != nil {
		return res.ErrInternalServer()
	}

	gender := new(entity.Gender)
	if data.Gender != "" {
		g := entity.Gender(data.Gender)
		gender = &g
	} else {
		gender = nil
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
		if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(data.OldPassword)); err != nil {
			return res.ErrForbidden(res.IncorrectOldPassword)
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return res.ErrInternalServer()
		}

		user.Password = string(hashedPassword)
	}

	if data.Photo != nil {
		contentType := data.Photo.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			return res.ErrUnprocessableEntity(res.PhotoOnly)
		}

		if data.Photo.Size > u.env.MaxUploadSize*1024*1024 {
			return res.ErrEntityTooLarge(res.PhotoSizeLimit)
		}

		src, err := data.Photo.Open()
		if err != nil {
			return res.ErrInternalServer()
		}

		defer src.Close()

		bucket := u.env.SupabaseBucket
		path := "profiles/" + uuid.NewString() + filepath.Ext(data.Photo.Filename)

		publicURL, err := u.Supabase.UploadFile(bucket, path, contentType, src)
		if err != nil {
			return res.ErrInternalServer()
		}

		user.PhotoURL = publicURL

		if userDB.PhotoURL != u.env.DefaultProfilePhotoURL {
			oldPhotoURL := userDB.PhotoURL
			index := strings.Index(oldPhotoURL, bucket)
			oldPhotoPath := oldPhotoURL[index+len(bucket+"/"):]

			if err = u.Supabase.DeleteFile(bucket, oldPhotoPath); err != nil {
				return res.ErrInternalServer()
			}
		}
	}

	if err := u.UserRepository.Update(user); err != nil {
		return res.ErrInternalServer()
	}

	return nil
}
