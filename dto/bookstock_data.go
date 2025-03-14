package dto

import (
	"time"

	"github.com/google/uuid"
)

type BookstockCreateRequest struct {
	Code   string    `json:"code" validate:"required"`
	BookID uuid.UUID `json:"book_id" validate:"required"`
}

type BookstockUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=AVAILABLE BORROWED DAMAGED LOST"`
}

type BookstockResponse struct {
	Code       string        `json:"code"`
	BookID     uuid.UUID     `json:"book_id"`
	Book       *BookResponse `json:"book,omitempty"`
	Status     string        `json:"status"`
	BorrowedID *uuid.UUID    `json:"borrowed_id"`
	BorrowedAt *time.Time    `json:"borrowed_at"`
}
