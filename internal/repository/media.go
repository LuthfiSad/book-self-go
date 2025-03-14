package repository

import (
	"go-rest-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaRepositoryImpl struct {
	db *gorm.DB
}

func NewMediaRepositoryImpl(db *gorm.DB) domain.MediaRepository {
	return &MediaRepositoryImpl{db: db}
}

func (r *MediaRepositoryImpl) FindAll() ([]domain.Media, error) {
	var media []domain.Media
	err := r.db.Find(&media).Error
	return media, err
}

func (r *MediaRepositoryImpl) FindByID(id uuid.UUID) (*domain.Media, error) {
	var media domain.Media
	if err := r.db.First(&media, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepositoryImpl) FindByFileName(fileName string) (*domain.Media, error) {
	var media domain.Media
	if err := r.db.First(&media, "path LIKE ?", "%"+fileName+"%").Error; err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepositoryImpl) Create(media *domain.Media) error {
	return r.db.Create(media).Error
}

func (r *MediaRepositoryImpl) Update(media *domain.Media) error {
	return r.db.Save(media).Error
}

func (r *MediaRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Media{}, id).Error
}
