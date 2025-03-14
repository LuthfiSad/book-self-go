package repository

import (
	"go-rest-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerRepositoryImpl(db *gorm.DB) domain.CustomerRepository {
	return &CustomerRepositoryImpl{db: db}
}

func (r *CustomerRepositoryImpl) FindAll() ([]domain.Customer, error) {
	var customers []domain.Customer
	err := r.db.Find(&customers).Error
	return customers, err
}

func (r *CustomerRepositoryImpl) FindByID(id uuid.UUID) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepositoryImpl) FindByCode(code string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.Where("code = ?", code).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepositoryImpl) Create(customer *domain.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepositoryImpl) Update(customer *domain.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Customer{}, id).Error
}
