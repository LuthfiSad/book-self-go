package domain

import (
	"go-rest-api/dto"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Path      string    `gorm:"size:255;not null" json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Books     []Book    `gorm:"foreignKey:CoverID" json:"books,omitempty"`
}

type MediaRepository interface {
	FindAll() ([]Media, error)
	FindByID(id uuid.UUID) (*Media, error)
	FindByFileName(fileName string) (*Media, error)
	Create(media *Media) error
	Update(media *Media) error
	Delete(id uuid.UUID) error
}

type MediaService interface {
	GetAllMedia() ([]dto.MediaResponse, error)
	GetMediaByID(id uuid.UUID) (*dto.MediaResponse, error)
	UploadMedia(file *multipart.FileHeader) (*dto.MediaResponse, error)
	DeleteMedia(id uuid.UUID) error
}
