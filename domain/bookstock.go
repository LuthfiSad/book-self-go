package domain

import (
	"go-rest-api/dto"
	"time"

	"github.com/google/uuid"
)

type BookStock struct {
	Code             string            `gorm:"primaryKey;size:50" json:"code"`
	BookID           uuid.UUID         `gorm:"not null" json:"book_id"`
	Book             Book              `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Status           string            `gorm:"size:50;not null" json:"status"` // Available, Borrowed, Damaged, Lost
	BorrowedID       *uuid.UUID        `json:"borrowed_id"`
	BorrowedAt       *time.Time        `json:"borrowed_at"`
	BookTransactions []BookTransaction `gorm:"foreignKey:StockCode;references:Code" json:"book_transactions,omitempty"`
}

type BookstockRepository interface {
	FindAll() ([]BookStock, error)
	FindByCode(code string) (*BookStock, error)
	FindByBookID(bookID uuid.UUID) ([]BookStock, error)
	FindAvailableByBookID(bookID uuid.UUID) ([]BookStock, error)
	Create(bookstock *BookStock) error
	Update(bookstock *BookStock) error
	Delete(code string) error
}

type BookstockService interface {
	GetAllBookstocks() ([]dto.BookstockResponse, error)
	GetBookstockByCode(code string) (*dto.BookstockResponse, error)
	GetBookstocksByBookID(bookID uuid.UUID) ([]dto.BookstockResponse, error)
	GetAvailableBookstocksByBookID(bookID uuid.UUID) ([]dto.BookstockResponse, error)
	CreateBookstock(req dto.BookstockCreateRequest) (*dto.BookstockResponse, error)
	UpdateBookstock(code string, req dto.BookstockUpdateRequest) (*dto.BookstockResponse, error)
	DeleteBookstock(code string) error
}
