package service

import (
	"context"
	"fmt"
	"strings"

	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/config"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type mediaService struct {
	mediaRepo   domain.MediaRepository
	bookService domain.BookService
	config      *config.Config
}

func NewMediaService(mediaRepo domain.MediaRepository, bookService domain.BookService, config *config.Config) domain.MediaService {
	return &mediaService{
		mediaRepo:   mediaRepo,
		bookService: bookService,
		config:      config,
	}
}

func (s *mediaService) GetAllMedia() ([]dto.MediaResponse, error) {
	media, err := s.mediaRepo.FindAll()
	if err != nil {
		return nil, err
	}

	mediaResponses := make([]dto.MediaResponse, 0, len(media))
	for _, m := range media {
		mediaResponses = append(mediaResponses, s.toMediaResponse(&m))
	}

	return mediaResponses, nil
}

func (s *mediaService) GetMediaByID(id uuid.UUID) (*dto.MediaResponse, error) {
	media, err := s.mediaRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := s.toMediaResponse(media)
	return &response, nil
}

func (s *mediaService) GetMediaByFileName(fileName string) (*dto.MediaResponse, error) {
	media, err := s.mediaRepo.FindByFileName(fileName)
	if err != nil {
		return nil, err
	}

	response := s.toMediaResponse(media)
	return &response, nil
}

func (s *mediaService) UploadMedia(newFileName string, filePath string) (*dto.MediaResponse, error) {
	media := &domain.Media{
		Path:      newFileName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.mediaRepo.Create(media); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to create media record: %w", err)
	}

	response := s.toMediaResponse(media)
	return &response, nil
}

func (s *mediaService) DeleteMedia(ctx context.Context, id uuid.UUID) error {
	media, err := s.mediaRepo.FindByID(id)
	if err != nil {
		return err
	}

	books, err := s.bookService.GetBookByCoverID(ctx, media.ID)
	if err != nil {
		return err
	}

	for _, book := range books {
		err := s.bookService.DeleteBookCover(ctx, book.ID)
		if err != nil {
			return fmt.Errorf("failed to delete book cover for book ID %s: %w", book.ID, err)
		}
	}

	errDelete := s.mediaRepo.Delete(id)
	if errDelete == nil {
		fileName := strings.TrimPrefix(media.Path, s.config.File.LinkCover+"/")
		filePath := filepath.Join(s.config.File.UploadPath, fileName)
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}

	return errDelete
}

func (s *mediaService) toMediaResponse(media *domain.Media) dto.MediaResponse {
	return dto.MediaResponse{
		ID:        media.ID,
		Path:      media.Path,
		CreatedAt: media.CreatedAt,
	}
}
