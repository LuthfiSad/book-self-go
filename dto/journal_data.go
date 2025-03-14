package dto

import (
	"time"

	"github.com/google/uuid"
)

type JournalCreateRequest struct {
	BookID     uuid.UUID `json:"book_id" validate:"required"`
	StockCode  string    `json:"stock_code" validate:"required"`
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	DueDate    string    `json:"due_date" validate:"required"` // Format: YYYY-MM-DD
}

type JournalReturnRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"` // User processing the return
}

type JournalResponse struct {
	ID         uuid.UUID          `json:"id"`
	BookID     uuid.UUID          `json:"book_id"`
	Book       *BookResponse      `json:"book,omitempty"`
	StockCode  string             `json:"stock_code"`
	BookStock  *BookstockResponse `json:"book_stock,omitempty"`
	CustomerID uuid.UUID          `json:"customer_id"`
	Customer   *CustomerResponse  `json:"customer,omitempty"`
	DueDate    time.Time          `json:"due_date"`
	Status     string             `json:"status"`
	BorrowedAt time.Time          `json:"borrowed_at"`
	ReturnAt   *time.Time         `json:"return_at"`
}
