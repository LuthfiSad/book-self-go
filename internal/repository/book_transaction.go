package repository

import (
	"go-rest-api/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookTransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewBookTransactionRepositoryImpl(db *gorm.DB) domain.BookTransactionRepository {
	return &BookTransactionRepositoryImpl{db: db}
}

func (r *BookTransactionRepositoryImpl) Find(filter map[string]interface{}) ([]domain.BookTransaction, error) {
	var transactions []domain.BookTransaction
	query := r.db.Preload("Book").Preload("Book.Cover").
		Preload("BookStock").Preload("Customer")

	// Apply filters dynamically
	for key, value := range filter {
		if key == "search" {
			query = query.
				Joins("JOIN customers ON customers.id = book_transactions.customer_id").
				Joins("JOIN books ON books.id = book_transactions.book_id").
				Where("customers.name LIKE ? OR book_transactions.title LIKE ? OR book_transactions.description LIKE ? OR books.title LIKE ?", "%"+value.(string)+"%", "%"+value.(string)+"%", "%"+value.(string)+"%", "%"+value.(string)+"%")
		} else {
			query = query.Where(key+" = ?", value)
		}
	}

	err := query.Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *BookTransactionRepositoryImpl) FindByID(id uuid.UUID) (*domain.BookTransaction, error) {
	var journal domain.BookTransaction
	err := r.db.Preload("Book").Preload("Book.Cover").
		Preload("BookStock").Preload("Customer").
		First(&journal, id).Error
	if err != nil {
		return nil, err
	}
	return &journal, nil
}

func (r *BookTransactionRepositoryImpl) Create(book_transaction *domain.BookTransaction) error {
	return r.db.Create(book_transaction).Error
}

func (r *BookTransactionRepositoryImpl) Update(book_transaction *domain.BookTransaction) error {
	return r.db.Save(book_transaction).Error
}

func (r *BookTransactionRepositoryImpl) UpdateBorrowStatus(id uuid.UUID, status string, borrowedAt *time.Time) error {
	return r.db.Model(&domain.BookTransaction{}).Where("id = ?", id).Updates(domain.BookTransaction{Status: status, BorrowedAt: borrowedAt}).Error
}

func (r *BookTransactionRepositoryImpl) UpdateReturnStatus(id uuid.UUID, status string, returnAt *time.Time) error {
	return r.db.Model(&domain.BookTransaction{}).Where("id = ?", id).Updates(domain.BookTransaction{Status: status, ReturnAt: returnAt}).Error
}

func (r *BookTransactionRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.BookTransaction{}, id).Error
}

func (r *BookTransactionRepositoryImpl) GetDB() *gorm.DB {
	return r.db
}
