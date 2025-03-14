package repository

import (
	"go-rest-api/domain"
	"go-rest-api/internal/constants"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type BookstockRepositoryImpl struct {
	db *gorm.DB
}

func NewBookstockRepositoryImpl(db *gorm.DB) domain.BookstockRepository {
	return &BookstockRepositoryImpl{db: db}
}

func (r *BookstockRepositoryImpl) FindAll() ([]domain.BookStock, error) {
	var bookstocks []domain.BookStock
	err := r.db.Preload("Book").Preload("Book.Cover").Find(&bookstocks).Error
	return bookstocks, err
}

func (r *BookstockRepositoryImpl) FindByCode(code string) (*domain.BookStock, error) {
	var bookstock domain.BookStock
	err := r.db.Preload("Book").Preload("Book.Cover").Where("code = ?", code).First(&bookstock).Error
	if err != nil {
		return nil, err
	}
	return &bookstock, nil
}

func (r *BookstockRepositoryImpl) FindByBookID(bookID uuid.UUID) ([]domain.BookStock, error) {
	var bookstocks []domain.BookStock
	err := r.db.Preload("Book").Preload("Book.Cover").Where("book_id = ?", bookID).Find(&bookstocks).Error
	return bookstocks, err
}

func (r *BookstockRepositoryImpl) FindAvailableByBookID(bookID uuid.UUID) ([]domain.BookStock, error) {
	var bookstocks []domain.BookStock
	err := r.db.Preload("Book").Preload("Book.Cover").
		Where("book_id = ? AND status = ?", bookID, constants.StatusAvailable).
		Find(&bookstocks).Error
	return bookstocks, err
}

func (r *BookstockRepositoryImpl) Create(bookstock *domain.BookStock) error {
	return r.db.Create(bookstock).Error
}

func (r *BookstockRepositoryImpl) Update(bookstock *domain.BookStock) error {
	return r.db.Save(bookstock).Error
}

func (r *BookstockRepositoryImpl) Delete(code string) error {
	return r.db.Delete(&domain.BookStock{}, "code = ?", code).Error
}
