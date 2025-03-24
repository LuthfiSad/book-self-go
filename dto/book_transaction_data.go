package dto

import (
	"time"

	"github.com/google/uuid"
)

type BookTransactionCreateRequest struct {
	BookID     uuid.UUID `json:"book_id" validate:"required"`
	StockCode  string    `json:"stock_code" validate:"required"`
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	DueDate    string    `json:"due_date" validate:"required"`
}

type BookTransactionUpdateRequest struct {
	ID         uuid.UUID `json:"id" validate:"omitempty"`
	BookID     uuid.UUID `json:"book_id" validate:"omitempty"`
	StockCode  string    `json:"stock_code" validate:"omitempty"`
	CustomerID uuid.UUID `json:"customer_id" validate:"omitempty"`
	DueDate    string    `json:"due_date" validate:"omitempty"`
	Status     string    `json:"status" validate:"omitempty"`
}

type BookTransactionUpdateStatusRequest struct {
	ID     uuid.UUID `json:"id" validate:"required"`
	Status string    `json:"status" validate:"required"`
	Date   time.Time `json:"date" validate:"required"`
}

type BookTransactionResponse struct {
	ID         uuid.UUID          `json:"id"`
	BookID     uuid.UUID          `json:"book_id"`
	Book       *BookResponse      `json:"book,omitempty"`
	StockCode  string             `json:"stock_code"`
	BookStock  *BookstockResponse `json:"book_stock,omitempty"`
	CustomerID uuid.UUID          `json:"customer_id"`
	Customer   *CustomerResponse  `json:"customer,omitempty"`
	DueDate    time.Time          `json:"due_date"`
	Status     string             `json:"status"`
	BorrowedAt *time.Time         `json:"borrowed_at"`
	ReturnAt   *time.Time         `json:"return_at"`
}
