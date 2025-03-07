package dto

import (
	"github.com/google/uuid"
	"mime/multipart"
	"time"
)

type GetRatingResponse struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	Star      int       `json:"star"`
	Feedback  string    `json:"feedback"`
	Photo     string    `json:"photo"`
	Date      time.Time `json:"date"`
}

type CreateRatingRequest struct {
	ProductID string                `form:"product_id" validate:"required,uuid"`
	Star      int                   `form:"star" validate:"required,min=1,max=5"`
	Feedback  string                `form:"feedback"`
	Photo     *multipart.FileHeader `form:"photo"`
}

type UpdateRatingRequest struct {
	Star     int                   `form:"star" validate:"required,min=1,max=5"`
	Feedback string                `form:"feedback"`
	Photo    *multipart.FileHeader `form:"photo"`
}

type RatingParam struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	UserID    uuid.UUID
}
