package dto

type UpdateUser struct {
	Name        string `json:"name"`
	OldPassword string `json:"old_password" validate:"required_with=NewPassword,omitempty,min=6"`
	NewPassword string `json:"new_password" validate:"omitempty,min=6"`
}
