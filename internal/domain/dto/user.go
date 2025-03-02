package dto

import "mime/multipart"

type UpdateUserRequest struct {
	Name        string                `form:"name"`
	Phone       string                `form:"phone"`
	Address     string                `form:"address"`
	BornDate    string                `form:"born_date" validate:"omitempty,datetime=2006-01-02"`
	Gender      string                `form:"gender" validate:"omitempty,oneof=male female"`
	Photo       *multipart.FileHeader `form:"photo" validate:"omitempty,image"`
	OldPassword string                `form:"old_password" validate:"required_with=NewPassword,omitempty,min=6"`
	NewPassword string                `form:"new_password" validate:"omitempty,min=6"`
}

type UpdateUserResponse struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	BornDate string `json:"born_date"`
	Gender   string `json:"gender"`
}

func (r *UpdateUserRequest) ToResponse() UpdateUserResponse {
	return UpdateUserResponse{
		Name:     r.Name,
		Phone:    r.Phone,
		Gender:   r.Gender,
		Address:  r.Address,
		BornDate: r.BornDate,
	}
}
