package domain

import (
	"go-rest-api/dto"
	"time"

	"github.com/google/uuid"
)

type Journal struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	BookID     uuid.UUID  `gorm:"not null" json:"book_id"`
	Book       Book       `gorm:"foreignKey:BookID" json:"book,omitempty"`
	StockCode  string     `gorm:"size:50;not null" json:"stock_code"`
	BookStock  BookStock  `gorm:"foreignKey:StockCode;references:Code" json:"book_stock,omitempty"`
	CustomerID uuid.UUID  `gorm:"not null" json:"customer_id"`
	Customer   Customer   `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	DueDate    time.Time  `json:"due_date"`
	Status     string     `gorm:"size:50;not null" json:"status"` // Borrowed, Returned, Overdue
	BorrowedAt time.Time  `json:"borrowed_at"`
	ReturnAt   *time.Time `json:"return_at"`
	Charges    []Charge   `gorm:"foreignKey:JournalID" json:"charges,omitempty"`
}

type JournalRepository interface {
	FindAll() ([]Journal, error)
	FindByID(id uuid.UUID) (*Journal, error)
	FindByCustomerID(customerID uuid.UUID) ([]Journal, error)
	FindByBookID(bookID uuid.UUID) ([]Journal, error)
	FindByStockCode(stockCode string) ([]Journal, error)
	FindActive() ([]Journal, error)
	FindOverdue() ([]Journal, error)
	Create(journal *Journal) error
	Update(journal *Journal) error
	UpdateStatus(id uuid.UUID, status string, returnAt *time.Time) error
	Delete(id uuid.UUID) error
}

type JournalService interface {
	GetAllJournals() ([]dto.JournalResponse, error)
	GetJournalByID(id uuid.UUID) (*dto.JournalResponse, error)
	GetJournalsByCustomerID(customerID uuid.UUID) ([]dto.JournalResponse, error)
	GetJournalsByBookID(bookID uuid.UUID) ([]dto.JournalResponse, error)
	GetActiveJournals() ([]dto.JournalResponse, error)
	GetOverdueJournals() ([]dto.JournalResponse, error)
	BorrowBook(req dto.JournalCreateRequest) (*dto.JournalResponse, error)
	ReturnBook(id uuid.UUID, req dto.JournalReturnRequest) (*dto.JournalResponse, error)
	DeleteJournal(id uuid.UUID) error
}
