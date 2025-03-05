package entity

import (
	"ambic/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rating struct {
	ID        uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	ProductID uuid.UUID `gorm:"type:varchar(36);not null;uniqueIndex:idx_product_user"`
	UserID    uuid.UUID `gorm:"type:varchar(36);not null;uniqueIndex:idx_product_user"`
	Star      int       `gorm:"type:int(8);not null"`
	Feedback  string    `gorm:"type:varchar(1000)"`
	PhotoURL  string    `gorm:"type:varchar(255)"`
}

func (r *Rating) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewUUID()
	r.ID = id
	return
}

func (r *Rating) ParseDTOGet() dto.GetRatingResponse {
	return dto.GetRatingResponse{
		ID:        r.ID,
		ProductID: r.ProductID.String(),
		UserID:    r.UserID.String(),
		Star:      r.Star,
		Feedback:  r.Feedback,
		Photo:     r.PhotoURL,
	}
}
