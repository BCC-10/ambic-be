package dto

import "time"

type UpdateUserRequest struct {
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	BornDate    time.Time `json:"born_date"`
	Gender      string    `json:"gender" validate:"omitempty,oneof=man woman"`
	PhotoURL    string    `json:"photo_url" validate:"omitempty,url"`
	OldPassword string    `json:"old_password" validate:"required_with=NewPassword,omitempty,min=6"`
	NewPassword string    `json:"new_password" validate:"omitempty,min=6"`
}
