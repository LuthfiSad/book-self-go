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

func (r *BookRepositoryImpl) FindBooks(ctx context.Context, page, perPage int, search string, cover_id *uuid.UUID) ([]domain.Book, int64, error) {
	var books []domain.Book
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Book{})

	if search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if cover_id != nil && *cover_id != uuid.Nil {
		query = query.Where("cover_id = ?", *cover_id)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && perPage > 0 {
		query = query.Offset((page - 1) * perPage).Limit(perPage)
	}

	err = query.Preload("Cover").Find(&books).Error

	return books, total, err
}

func (r *BookRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	var book domain.Book
	err := r.db.WithContext(ctx).Preload("Cover").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
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
