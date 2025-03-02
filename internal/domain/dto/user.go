package dto

import "mime/multipart"

type UpdateUserRequest struct {
	Name        string                `json:"name" form:"name"`
	Phone       string                `json:"phone" form:"phone"`
	Address     string                `json:"address" form:"address"`
	BornDate    string                `json:"born_date" form:"born_date" validate:"omitempty,datetime=2006-01-02"`
	Gender      string                `json:"gender" form:"gender" validate:"omitempty,oneof=male female"`
	Photo       *multipart.FileHeader `json:"-" form:"photo" validate:"omitempty,image"`
	OldPassword string                `json:"old_password" form:"old_password" validate:"required_with=NewPassword,omitempty,min=6"`
	NewPassword string                `json:"new_password" form:"new_password" validate:"omitempty,min=6"`
}

type UpdateUserResponse struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	BornDate string `json:"born_date"`
	Gender   string `json:"gender"`
	PhotoURL string `json:"photo"`
}

func (r *UpdateUserRequest) ToResponse(val ...string) UpdateUserResponse {
	res := UpdateUserResponse{
		Name:     r.Name,
		Phone:    r.Phone,
		Gender:   r.Gender,
		Address:  r.Address,
		BornDate: r.BornDate,
	}

	if len(val) == 1 {
		res.PhotoURL = val[0]
	}

	return res
}
