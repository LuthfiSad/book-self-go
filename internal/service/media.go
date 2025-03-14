package service

import (
	"errors"
	"fmt"
	"strings"

	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/config"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type mediaService struct {
	mediaRepo domain.MediaRepository
	config    *config.Config
}

func NewMediaService(mediaRepo domain.MediaRepository, config *config.Config) domain.MediaService {
	return &mediaService{
		mediaRepo: mediaRepo,
		config:    config,
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

func (s *mediaService) UploadMedia(file *multipart.FileHeader) (*dto.MediaResponse, error) {
	maxSize, err := strconv.ParseInt(s.config.File.MaxUploadSize, 10, 64)
	if err != nil {
		maxSize = 5
	}
	maxSize = maxSize * 1024 * 1024

	if file.Size > maxSize {
		return nil, errors.New("file size exceeds maximum upload size")
	}

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return nil, errors.New("unsupported file type")
	}

	uploadPath := s.config.File.UploadPath
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadPath, newFileName)

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	media := &domain.Media{
		Path:      s.config.File.LinkCover + "/" + newFileName,
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

func (s *mediaService) DeleteMedia(id uuid.UUID) error {
	media, err := s.mediaRepo.FindByID(id)
	if err != nil {
		return err
	}

	filePath := filepath.Join(s.config.File.UploadPath, media.Path)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return s.mediaRepo.Delete(id)
}

func (s *mediaService) toMediaResponse(media *domain.Media) dto.MediaResponse {
	return dto.MediaResponse{
		ID:        media.ID,
		Path:      media.Path,
		CreatedAt: media.CreatedAt,
	}
}
