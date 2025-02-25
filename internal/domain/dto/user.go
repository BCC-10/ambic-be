package dto

import "github.com/google/uuid"

type Register struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RequestOTP struct {
	Email string `json:"email" validate:"required,email"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserParam struct {
	Id    uuid.UUID
	Email string
}
