package repository

import (
	"context"
	"go-rest-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookRepositoryImpl struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) domain.BookRepository {
	return &BookRepositoryImpl{db: db}
}

func (r *BookRepositoryImpl) FindAll(ctx context.Context) ([]domain.Book, error) {
	var books []domain.Book
	err := r.db.WithContext(ctx).Preload("Cover").Find(&books).Error
	return books, err
}

func (r *BookRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	var book domain.Book
	err := r.db.WithContext(ctx).Preload("Cover").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepositoryImpl) FindByCoverID(ctx context.Context, id uuid.UUID) ([]domain.Book, error) {
	var books []domain.Book
	err := r.db.WithContext(ctx).Preload("Cover").Where("cover_id = ?", id).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepositoryImpl) Create(ctx context.Context, book *domain.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *BookRepositoryImpl) Update(ctx context.Context, book *domain.Book) error {
	return r.db.WithContext(ctx).Save(book).Error
}

func (r *BookRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Book{}, id).Error
}
