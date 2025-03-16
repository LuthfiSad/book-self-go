package service

import (
	"context"
	"errors"
	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/config"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type bookService struct {
	bookRepo  domain.BookRepository
	mediaRepo domain.MediaRepository
	config    *config.Config
}

func NewBookService(bookRepo domain.BookRepository, mediaRepo domain.MediaRepository, config *config.Config) domain.BookService {
	return &bookService{
		bookRepo:  bookRepo,
		mediaRepo: mediaRepo,
		config:    config,
	}
}

func (s *bookService) GetAllBooks(ctx context.Context) ([]dto.BookResponse, error) {
	books, err := s.bookRepo.FindAll(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	bookResponses := make([]dto.BookResponse, 0, len(books))
	for _, book := range books {
		bookResponses = append(bookResponses, s.toBookResponse(&book))
	}

	return bookResponses, nil
}

func (s *bookService) GetBookByID(ctx context.Context, id uuid.UUID) (*dto.BookResponse, error) {
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	response := s.toBookResponse(book)
	return &response, nil
}

func (s *bookService) GetBookByCoverID(ctx context.Context, id uuid.UUID) ([]dto.BookResponse, error) {
	books, err := s.bookRepo.FindByCoverID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	bookResponses := make([]dto.BookResponse, 0, len(books))
	for _, book := range books {
		bookResponses = append(bookResponses, s.toBookResponse(&book))
	}

	return bookResponses, nil
}

func (s *bookService) CreateBook(ctx context.Context, req dto.BookCreateRequest) (*dto.BookResponse, error) {
	book := &domain.Book{
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.CoverID != nil {
		media, err := s.mediaRepo.FindByID(*req.CoverID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return nil, errors.New("invalid cover ID: media not found")
		}

		book.CoverID = &media.ID
	}

	if err := s.bookRepo.Create(ctx, book); err != nil {
		return nil, err
	}

	response := s.toBookResponse(book)
	return &response, nil
}

func (s *bookService) UpdateBook(ctx context.Context, id uuid.UUID, req dto.BookUpdateRequest) (*dto.BookResponse, error) {
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		book.Title = req.Title
	}

	if req.Description != "" {
		book.Description = req.Description
	}

	if req.CoverID != nil {
		media, err := s.mediaRepo.FindByID(*req.CoverID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return nil, errors.New("invalid cover ID: media not found")
		}

		book.CoverID = &media.ID
		book.Cover = media
	}

	book.UpdatedAt = time.Now()

	if err := s.bookRepo.Update(ctx, book); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	response := s.toBookResponse(book)
	return &response, nil
}

func (s *bookService) DeleteBook(ctx context.Context, id uuid.UUID) error {
	if _, err := s.bookRepo.FindByID(ctx, id); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	return s.bookRepo.Delete(ctx, id)
}

func (s *bookService) DeleteBookCover(ctx context.Context, id uuid.UUID) error {
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	if book.Cover == nil && book.CoverID == nil {
		return errors.New("book cover not found")
	}

	book.Cover = nil
	book.CoverID = nil
	book.UpdatedAt = time.Now()

	return s.bookRepo.Update(ctx, book)
}

func (s *bookService) toBookResponse(book *domain.Book) dto.BookResponse {
	response := dto.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}

	if book.Cover != nil {
		response.Cover = &dto.MediaResponse{
			ID:        book.Cover.ID,
			Path:      book.Cover.Path,
			CreatedAt: book.Cover.CreatedAt,
		}
	}

	return response
}
