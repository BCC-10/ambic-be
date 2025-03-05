package dto

import "mime/multipart"

type GetRatingResponse struct {
	ID        uint   `json:"id"`
	ProductID string `json:"product_id"`
	UserID    string `json:"user_id"`
	Star      int    `json:"star"`
	Feedback  string `json:"feedback"`
	Photo     string `json:"photo"`
}

type CreateRatingRequest struct {
	ProductID string                `form:"product_id" validate:"required,uuid"`
	UserID    string                `form:"user_id" validate:"required,uuid"`
	Star      int                   `form:"star" validate:"required,min=1,max=5"`
	Feedback  string                `form:"feedback"`
	Photo     *multipart.FileHeader `form:"photo"`
}

type UpdateRatingRequest struct {
	ProductID string                `form:"product_id" validate:"required,uuid"`
	UserID    string                `form:"user_id" validate:"required,uuid"`
	Star      int                   `form:"star" validate:"required,min=1,max=5"`
	Feedback  string                `form:"feedback"`
	Photo     *multipart.FileHeader `form:"photo"`
}

type RatingParam struct {
	ID        uint
	ProductID string
	UserID    string
}
