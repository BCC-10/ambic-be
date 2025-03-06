package dto

import "github.com/google/uuid"

type GetBusinessTypeResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type BusinessTypeParam struct {
	ID uuid.UUID
}
