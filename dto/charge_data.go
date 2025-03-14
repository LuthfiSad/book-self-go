package dto

import (
	"time"

	"github.com/google/uuid"
)

type ChargeCreateRequest struct {
	JournalID    uuid.UUID `json:"journal_id" validate:"required"`
	DaysLate     int       `json:"days_late" validate:"required,min=1"`
	DailyLateFee float64   `json:"daily_late_fee" validate:"required,min=0"`
	UserID       uuid.UUID `json:"user_id" validate:"required"` // User processing the charge
}

type ChargeUpdateRequest struct {
	DaysLate     int     `json:"days_late" validate:"required,min=1"`
	DailyLateFee float64 `json:"daily_late_fee" validate:"required,min=0"`
}

type ChargeResponse struct {
	ID           uuid.UUID        `json:"id"`
	JournalID    uuid.UUID        `json:"journal_id"`
	Journal      *JournalResponse `json:"journal,omitempty"`
	DaysLate     int              `json:"days_late"`
	DailyLateFee float64          `json:"daily_late_fee"`
	Total        float64          `json:"total"`
	UserID       uuid.UUID        `json:"user_id"`
	User         *UserData        `json:"user,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
}
