package dto

import "github.com/google/uuid"

type GetNotificationResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Link     string `json:"link"`
	Photo    string `json:"photo"`
	Datetime string `json:"datetime"`
}

type NotificationParam struct {
	UserID uuid.UUID
}
