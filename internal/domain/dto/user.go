package dto

import "mime/multipart"

type GetUserResponse struct {
	ID       string             `json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Name     string             `json:"name"`
	Phone    string             `json:"phone"`
	Address  string             `json:"address"`
	Gender   string             `json:"gender"`
	Photo    string             `json:"photo"`
	Partner  GetPartnerResponse `json:"partner"`
}

type UpdateUserRequest struct {
	Name        string                `form:"name"`
	Phone       string                `form:"phone"`
	Address     string                `form:"address"`
	Gender      string                `form:"gender" validate:"omitempty,oneof=male female"`
	Photo       *multipart.FileHeader `form:"photo"`
	OldPassword string                `form:"old_password" validate:"required_with=NewPassword,omitempty,min=6"`
	NewPassword string                `form:"new_password" validate:"omitempty,min=6"`
}
