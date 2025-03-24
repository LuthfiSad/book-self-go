package domain

import (
	"context"
	"go-rest-api/dto"
	"time"

	"github.com/google/uuid"
)

type BookTransaction struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	BookID     uuid.UUID  `gorm:"not null" json:"book_id"`
	Book       Book       `gorm:"foreignKey:BookID" json:"book,omitempty"`
	StockCode  string     `gorm:"size:50;not null" json:"stock_code"`
	BookStock  BookStock  `gorm:"foreignKey:StockCode;references:Code" json:"book_stock,omitempty"`
	CustomerID uuid.UUID  `gorm:"not null" json:"customer_id"`
	Customer   Customer   `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	DueDate    time.Time  `json:"due_date"`
	Status     string     `gorm:"size:50;not null" json:"status"` // Borrowed, Returned, Overdue
	BorrowedAt *time.Time `json:"borrowed_at"`
	ReturnAt   *time.Time `json:"return_at"`
	Charges    []Charge   `gorm:"foreignKey:BookTransactionID" json:"charges,omitempty"`
}

type BookTransactionRepository interface {
	Find(filter map[string]interface{}) ([]BookTransaction, error)
	FindByID(id uuid.UUID) (*BookTransaction, error)
	Create(book_transaction *BookTransaction) error
	Update(book_transaction *BookTransaction) error
	UpdateBorrowStatus(id uuid.UUID, status string, borrowedAt *time.Time) error
	UpdateReturnStatus(id uuid.UUID, status string, returnAt *time.Time) error
	Delete(id uuid.UUID) error
}

type BookTransactionService interface {
	GetAllBookTransactions(filter map[string]interface{}) ([]dto.BookTransactionResponse, error)
	CreateBookTransaction(ctx context.Context, req dto.BookTransactionCreateRequest) (*dto.BookTransactionResponse, error)
	UpdateBookTransaction(ctx context.Context, id uuid.UUID, req dto.BookTransactionUpdateRequest) (*dto.BookTransactionResponse, error)
	ReturnBookTransaction(ctx context.Context, req dto.BookTransactionUpdateStatusRequest) (*dto.BookTransactionResponse, error)
	DeleteBookTransaction(ctx context.Context, id uuid.UUID) error
}
