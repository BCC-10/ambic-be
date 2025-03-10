package entity

import (
	"ambic/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `gorm:"type:char(36)"`
	Title     string    `gorm:"type:varchar(255); not null"`
	Content   string    `gorm:"type:varchar(255);not null"`
	Link      string    `gorm:"type:varchar(255)"`
	Button    string    `gorm:"type:varchar(255)"`
	PhotoURL  string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	n.ID, _ = uuid.NewUUID()
	return
}

func (n *Notification) ParseDTOGet() dto.GetNotificationResponse {
	return dto.GetNotificationResponse{
		ID:       n.ID.String(),
		Title:    n.Title,
		Content:  n.Content,
		Link:     n.Link,
		Button:   n.Button,
		Photo:    n.PhotoURL,
		Datetime: n.CreatedAt.String(),
	}
}
