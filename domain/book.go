package domain

import (
	"context"
	"go-rest-api/dto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID               uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Title            string            `gorm:"size:255;not null" json:"title"`
	Description      string            `gorm:"type:text" json:"description"`
	CoverID          *uuid.UUID        `json:"cover_id"`
	Cover            *Media            `gorm:"foreignKey:CoverID" json:"cover,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `gorm:"index" json:"-"`
	BookStocks       []BookStock       `gorm:"foreignKey:BookID" json:"book_stocks,omitempty"`
	BookTransactions []BookTransaction `gorm:"foreignKey:BookID" json:"book_transactions,omitempty"`
}

type BookRepository interface {
	FindBooks(ctx context.Context, page, perPage int, search string, cover_id *uuid.UUID) ([]Book, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Book, error)
	Create(ctx context.Context, book *Book) error
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BookService interface {
	GetBooks(ctx context.Context, page, perPage int, search string, cover_id *uuid.UUID) (*dto.PaginatedResponseData[[]dto.BookResponse], error)
	GetBookByID(ctx context.Context, id uuid.UUID) (*dto.BookResponse, error)
	CreateBook(ctx context.Context, req dto.BookCreateRequest) (*dto.BookResponse, error)
	UpdateBook(ctx context.Context, id uuid.UUID, req dto.BookUpdateRequest) (*dto.BookResponse, error)
	DeleteBook(ctx context.Context, id uuid.UUID) error
	DeleteBookCover(ctx context.Context, id uuid.UUID) error
}
