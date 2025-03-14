package dto

import (
	"time"

	"github.com/google/uuid"
)

type MediaResponse struct {
	ID        uuid.UUID `json:"id"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}
