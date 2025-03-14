package repository

import (
	"go-rest-api/domain"
	"go-rest-api/internal/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JournalRepositoryImpl struct {
	db *gorm.DB
}

func NewJournalRepositoryImpl(db *gorm.DB) domain.JournalRepository {
	return &JournalRepositoryImpl{db: db}
}

func (r *JournalRepositoryImpl) FindAll() ([]domain.Journal, error) {
	var journals []domain.Journal
	err := r.db.Preload("Book").Preload("Book.Cover").
		Preload("BookStock").Preload("Customer").
		Find(&journals).Error
	return journals, err
}

func (r *JournalRepositoryImpl) FindByID(id uuid.UUID) (*domain.Journal, error) {
	var journal domain.Journal
	err := r.db.Preload("Book").Preload("Book.Cover").
		Preload("BookStock").Preload("Customer").
		First(&journal, id).Error
	if err != nil {
		return nil, err
	}
	return &journal, nil
}

func (r *JournalRepositoryImpl) FindByCustomerID(customerID uuid.UUID) ([]domain.Journal, error) {
	var journals []domain.Journal
	err := r.db.Preload("Book").Preload("Book.Cover").
		Preload("BookStock").Preload("Customer").
		Where("customer_id = ?", customerID).
		Find(&journals).Error
	return journals, err
}

func (r *JournalRepositoryImpl) FindByBookID(bookID uuid.UUID) ([]domain.Journal, error) {
	var journals []domain.Journal
	err := r.db.Preload("Book").Preload("Book.Cover").
		Preload("BookStock").Preload("Customer").
		Where("book_id = ?", bookID).
		Find(&journals).Error
	return journals, err
}

func (r *JournalRepositoryImpl) FindByStockCode(code string) ([]domain.Journal, error) {
	var journals []domain.Journal
	err := r.db.Preload("Book").Preload("Book.Cover").
		Preload("BookStock").Preload("Customer").
		Where("stock_code = ?", code).
		Find(&journals).Error
	return journals, err
}

func (r *JournalRepositoryImpl) FindActive() ([]domain.Journal, error) {
	var journals []domain.Journal
	err := r.db.Preload("Book").Preload("Book.Cover").
		Preload("BookStock").Preload("Customer").
		Where("status = ?", constants.StatusActive).
		Find(&journals).Error
	return journals, err
}

func (r *JournalRepositoryImpl) FindOverdue() ([]domain.Journal, error) {
	var journals []domain.Journal
	err := r.db.Preload("Book").Preload("Book.Cover").
		Preload("BookStock").Preload("Customer").
		Where("status = ? AND return_at < ?", constants.StatusActive, time.Now()).
		Find(&journals).Error
	return journals, err
}

func (r *JournalRepositoryImpl) FindByCode(code string) (*domain.Journal, error) {
	var journal domain.Journal
	err := r.db.Preload("Book").Preload("Book.Cover").
		Preload("BookStock").Preload("Customer").
		Where("code = ?", code).
		First(&journal).Error
	if err != nil {
		return nil, err
	}
	return &journal, nil
}

func (r *JournalRepositoryImpl) Create(journal *domain.Journal) error {
	return r.db.Create(journal).Error
}

func (r *JournalRepositoryImpl) Update(journal *domain.Journal) error {
	return r.db.Save(journal).Error
}

func (r *JournalRepositoryImpl) UpdateStatus(id uuid.UUID, status string, returnAt *time.Time) error {
	return r.db.Model(&domain.Journal{}).Where("id = ?", id).Updates(domain.Journal{Status: status, ReturnAt: returnAt}).Error
}

func (r *JournalRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Journal{}, id).Error
}
