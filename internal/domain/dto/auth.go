package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Name       string `json:"name" validate:"required"`
	Username   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	IsVerified bool   `json:"is_verified"`
}

type RegisterResponse struct {
	Name     string `json:"name"`
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

type LoginResponse struct {
	Identifier string `json:"identifier"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Token    string `json:"token" validate:"required"`
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
		Name:     r.Name,
		Username: r.Username,
		Email:    r.Email,
	}
}

func (r LoginRequest) AsResponse() LoginResponse {
	return LoginResponse{
		Identifier: r.Identifier,
	}
}
