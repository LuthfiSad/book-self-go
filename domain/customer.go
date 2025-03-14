package domain

import (
	"go-rest-api/dto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Code      string         `gorm:"size:50;not null;unique" json:"code"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Journals  []Journal      `gorm:"foreignKey:CustomerID" json:"journals,omitempty"`
}

type CustomerRepository interface {
	FindAll() ([]Customer, error)
	FindByID(id uuid.UUID) (*Customer, error)
	FindByCode(code string) (*Customer, error)
	Create(customer *Customer) error
	Update(customer *Customer) error
	Delete(id uuid.UUID) error
}

type CustomerService interface {
	GetAllCustomers() ([]dto.CustomerResponse, error)
	GetCustomerByID(id uuid.UUID) (*dto.CustomerResponse, error)
	GetCustomerByCode(code string) (*dto.CustomerResponse, error)
	CreateCustomer(req dto.CustomerCreateRequest) (*dto.CustomerResponse, error)
	UpdateCustomer(id uuid.UUID, req dto.CustomerUpdateRequest) (*dto.CustomerResponse, error)
	DeleteCustomer(id uuid.UUID) error
}
