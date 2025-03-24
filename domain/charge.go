package domain

import (
	"go-rest-api/dto"
	"time"

	"github.com/google/uuid"
)

type Charge struct {
	ID                uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	BookTransactionID uuid.UUID       `gorm:"not null" json:"book_transaction_id"`
	BookTransaction   BookTransaction `gorm:"foreignKey:BookTransactionID" json:"book_transaction,omitempty"`
	DaysLate          int             `gorm:"not null" json:"days_late"`
	DailyLateFee      float64         `gorm:"not null" json:"daily_late_fee"`
	Total             float64         `gorm:"not null" json:"total"`
	UserID            uuid.UUID       `gorm:"not null" json:"user_id"`
	User              User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
}

type ChargeRepository interface {
	FindAll() ([]Charge, error)
	FindByID(id uuid.UUID) (*Charge, error)
	FindByBookTransactionID(book_transactionID uuid.UUID) ([]Charge, error)
	Create(charge *Charge) error
	Update(charge *Charge) error
	Delete(id uuid.UUID) error
}

type ChargeService interface {
	GetAllCharges() ([]dto.ChargeResponse, error)
	GetChargeByID(id uuid.UUID) (*dto.ChargeResponse, error)
	GetChargesByBookTransactionID(book_transactionID uuid.UUID) ([]dto.ChargeResponse, error)
	CreateCharge(req dto.ChargeCreateRequest) (*dto.ChargeResponse, error)
	UpdateCharge(id uuid.UUID, req dto.ChargeUpdateRequest) (*dto.ChargeResponse, error)
	DeleteCharge(id uuid.UUID) error
}
