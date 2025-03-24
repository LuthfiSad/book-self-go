package repository

import (
	"go-rest-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChargeRepositoryImpl struct {
	db *gorm.DB
}

func NewChargeRepositoryImpl(db *gorm.DB) domain.ChargeRepository {
	return &ChargeRepositoryImpl{db: db}
}

func (r *ChargeRepositoryImpl) FindAll() ([]domain.Charge, error) {
	var charges []domain.Charge
	err := r.db.Preload("BookTransaction").Preload("BookTransaction.Book").
		Preload("BookTransaction.Customer").Preload("User").
		Find(&charges).Error
	return charges, err
}

func (r *ChargeRepositoryImpl) FindByID(id uuid.UUID) (*domain.Charge, error) {
	var charge domain.Charge
	err := r.db.Preload("BookTransaction").Preload("BookTransaction.Book").
		Preload("BookTransaction.Customer").Preload("User").
		First(&charge, id).Error
	if err != nil {
		return nil, err
	}
	return &charge, nil
}

func (r *ChargeRepositoryImpl) FindByBookTransactionID(book_transactionID uuid.UUID) ([]domain.Charge, error) {
	var charges []domain.Charge
	err := r.db.Preload("BookTransaction").Preload("BookTransaction.Book").
		Preload("BookTransaction.Customer").Preload("User").
		Where("book_transaction_id = ?", book_transactionID).
		Find(&charges).Error
	return charges, err
}

func (r *ChargeRepositoryImpl) FindByCustomerID(customerID uuid.UUID) ([]domain.Charge, error) {
	var charges []domain.Charge
	err := r.db.Preload("BookTransaction").Preload("BookTransaction.Book").
		Preload("BookTransaction.Customer").Preload("User").
		Joins("JOIN book_transactions ON charges.book_transaction_id = book_transactions.id").
		Where("book_transactions.customer_id = ?", customerID).
		Find(&charges).Error
	return charges, err
}

func (r *ChargeRepositoryImpl) Create(charge *domain.Charge) error {
	return r.db.Create(charge).Error
}

func (r *ChargeRepositoryImpl) Update(charge *domain.Charge) error {
	return r.db.Save(charge).Error
}

func (r *ChargeRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Charge{}, id).Error
}
