package dto

import "github.com/google/uuid"

type GetBusinessTypeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BusinessTypeParam struct {
	ID uuid.UUID
}
