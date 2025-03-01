package dto

type UpdateUserRequest struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	BornDate    string `json:"born_date" validate:"omitempty,datetime=2006-01-02"`
	Gender      string `json:"gender" validate:"omitempty,oneof=male female"`
	PhotoURL    string `json:"photo_url" validate:"omitempty,url"`
	OldPassword string `json:"old_password" validate:"required_with=NewPassword,omitempty,min=6"`
	NewPassword string `json:"new_password" validate:"omitempty,min=6"`
}

type UpdateUserResponse struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	BornDate string `json:"born_date" validate:"omitempty,datetime=2006-01-02"`
	Gender   string `json:"gender" validate:"omitempty,oneof=man woman"`
	PhotoURL string `json:"photo_url" validate:"omitempty,url"`
}

func (r *UpdateUserRequest) ToResponse() UpdateUserResponse {
	return UpdateUserResponse{
		Name:     r.Name,
		Phone:    r.Phone,
		Gender:   r.Gender,
		Address:  r.Address,
		BornDate: r.BornDate,
		PhotoURL: r.PhotoURL,
	}
}
