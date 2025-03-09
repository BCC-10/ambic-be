package dto

import (
	"github.com/google/uuid"
	"mime/multipart"
	"time"
)

type GetRatingRequest struct {
	ID        string `query:"id" validate:"omitempty,uuid"`
	ProductID string `query:"product_id" validate:"omitempty,uuid"`
	UserID    string `query:"user_id" validate:"omitempty,uuid"`
}

type ShowRatingRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

type GetRatingResponse struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	Star      int       `json:"star"`
	Feedback  string    `json:"feedback"`
	Photo     string    `json:"photo"`
	Datetime  time.Time `json:"datetime"`
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

type UpdateRatingParam struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

type RatingParam struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	UserID    uuid.UUID
}

func (g GetRatingRequest) ParseParam() RatingParam {
	param := new(RatingParam)

	if g.ID != "" {
		param.ID, _ = uuid.Parse(g.ID)
	}

	if g.ProductID != "" {
		param.ProductID, _ = uuid.Parse(g.ProductID)
	}

	if g.UserID != "" {
		param.UserID, _ = uuid.Parse(g.UserID)
	}

	return *param
}
