package dto

import (
	"time"

	"github.com/google/uuid"
)

type CustomerCreateRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type CustomerUpdateRequest struct {
	Name string `json:"name" validate:"required"`
}

type CustomerResponse struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
