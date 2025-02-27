package dto

import "github.com/google/uuid"

type Register struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RequestOTP struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyOTP struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required"`
}

type Login struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=6"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type UserParam struct {
	Id       uuid.UUID
	Email    string
	Username string
}

func (r Register) AsResponse() RegisterResponse {
	return RegisterResponse{
		Name:     r.Name,
		Username: r.Username,
		Email:    r.Email,
	}
}
