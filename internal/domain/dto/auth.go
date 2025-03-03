package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Username   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	IsVerified bool   `json:"is_verified"`
}

type RegisterResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RequestTokenRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=6"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Token    string `json:"token" validate:"required"`
}

type GoogleCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
	Error string `json:"error"`
}

type GoogleUserProfileResponse struct {
	Email      string
	Username   string
	Name       string
	IsVerified bool
}

type UserParam struct {
	Id       uuid.UUID
	Email    string
	Username string
}

type Empty struct {
}

func (r RegisterRequest) AsResponse() RegisterResponse {
	return RegisterResponse{
		Username: r.Username,
		Email:    r.Email,
	}
}
