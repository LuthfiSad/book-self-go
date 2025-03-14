package dto

import (
	"time"

	"github.com/google/uuid"
)

type BookCreateRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	CoverID     *uuid.UUID `json:"cover_id" validate:"omitempty"`
}

type BookUpdateRequest struct {
	Title       string     `json:"title" validate:"omitempty"`
	Description string     `json:"description" validate:"omitempty"`
	CoverID     *uuid.UUID `json:"cover_id" validate:"omitempty"`
}

type BookResponse struct {
	ID          uuid.UUID      `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	CoverID     *uuid.UUID     `json:"-"`
	Cover       *MediaResponse `json:"cover,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
