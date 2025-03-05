package entity

import (
	"ambic/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rating struct {
	ID        uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	ProductID uuid.UUID `gorm:"type:varchar(36);not null"`
	UserID    uuid.UUID `gorm:"type:varchar(36);not null"`
	Star      int       `gorm:"type:int(8);not null"`
	Feedback  string    `gorm:"type:varchar(1000);not null"`
	PhotoURL  string    `gorm:"type:varchar(255)"`
}

func (r *Rating) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Exec("ALTER TABLE ratings ADD CONSTRAINT unique_user_product UNIQUE (user_id, product_id)").Error
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
