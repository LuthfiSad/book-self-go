package service

import (
	"context"
	"errors"
	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/constants"
	"time"

	"github.com/google/uuid"
)

type bookstockService struct {
	bookstockRepo domain.BookstockRepository
	bookRepo      domain.BookRepository
}

func NewBookstockService(bookstockRepo domain.BookstockRepository, bookRepo domain.BookRepository) domain.BookstockService {
	return &bookstockService{
		bookstockRepo: bookstockRepo,
		bookRepo:      bookRepo,
	}
}

func (s *bookstockService) GetAllBookstocks() ([]dto.BookstockResponse, error) {
	bookstocks, err := s.bookstockRepo.FindAll()
	if err != nil {
		return nil, err
	}

	bookstockResponses := make([]dto.BookstockResponse, 0, len(bookstocks))
	for _, bookstock := range bookstocks {
		bookstockResponses = append(bookstockResponses, s.toBookstockResponse(&bookstock))
	}

	return bookstockResponses, nil
}

func (s *bookstockService) GetBookstockByCode(code string) (*dto.BookstockResponse, error) {
	bookstock, err := s.bookstockRepo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	response := s.toBookstockResponse(bookstock)
	return &response, nil
}

func (s *bookstockService) GetBookstocksByBookID(bookID uuid.UUID) ([]dto.BookstockResponse, error) {
	bookstocks, err := s.bookstockRepo.FindByBookID(bookID)
	if err != nil {
		return nil, err
	}

	bookstockResponses := make([]dto.BookstockResponse, 0, len(bookstocks))
	for _, bookstock := range bookstocks {
		bookstockResponses = append(bookstockResponses, s.toBookstockResponse(&bookstock))
	}

	return bookstockResponses, nil
}

func (s *bookstockService) GetAvailableBookstocksByBookID(bookID uuid.UUID) ([]dto.BookstockResponse, error) {
	bookstocks, err := s.bookstockRepo.FindAvailableByBookID(bookID)
	if err != nil {
		return nil, err
	}

	bookstockResponses := make([]dto.BookstockResponse, 0, len(bookstocks))
	for _, bookstock := range bookstocks {
		bookstockResponses = append(bookstockResponses, s.toBookstockResponse(&bookstock))
	}

	return bookstockResponses, nil
}

func (s *bookstockService) CreateBookstock(req dto.BookstockCreateRequest) (*dto.BookstockResponse, error) {
	book, err := s.bookRepo.FindByID(context.Background(), req.BookID)
	if err != nil {
		return nil, errors.New("invalid book ID: book not found")
	}

	existingBookstock, err := s.bookstockRepo.FindByCode(req.Code)
	if err == nil && existingBookstock != nil {
		return nil, errors.New("bookstock with this code already exists")
	}

	bookstock := &domain.BookStock{
		Code:   req.Code,
		BookID: req.BookID,
		Book:   *book,
		Status: constants.BookStockStatusAvailable, // Default status
	}

	if err := s.bookstockRepo.Create(bookstock); err != nil {
		return nil, err
	}

	response := s.toBookstockResponse(bookstock)
	return &response, nil
}

func (s *bookstockService) UpdateBookstock(code string, req dto.BookstockUpdateRequest) (*dto.BookstockResponse, error) {
	bookstock, err := s.bookstockRepo.FindByCode(code)
	if err != nil {
		return nil, errors.New("bookstock not found")
	}

	bookstock.Status = req.Status

	if req.Status == constants.BookStockStatusBorrowed && bookstock.BorrowedID == nil {
		dummyBorrowID := uuid.New()
		now := time.Now()
		bookstock.BorrowedID = &dummyBorrowID
		bookstock.BorrowedAt = &now
	} else if req.Status != constants.BookStockStatusBorrowed {
		bookstock.BorrowedID = nil
		bookstock.BorrowedAt = nil
	}

	if err := s.bookstockRepo.Update(bookstock); err != nil {
		return nil, err
	}

	response := s.toBookstockResponse(bookstock)
	return &response, nil
}

func (s *bookstockService) DeleteBookstock(code string) error {
	// Check if bookstock exists
	if _, err := s.bookstockRepo.FindByCode(code); err != nil {
		return errors.New("bookstock not found")
	}

	return s.bookstockRepo.Delete(code)
}

func (s *bookstockService) toBookstockResponse(bookstock *domain.BookStock) dto.BookstockResponse {
	response := dto.BookstockResponse{
		Code:       bookstock.Code,
		BookID:     bookstock.BookID,
		Status:     bookstock.Status,
		BorrowedID: bookstock.BorrowedID,
		BorrowedAt: bookstock.BorrowedAt,
	}

	if bookstock.Book.ID != uuid.Nil {
		bookResponse := &dto.BookResponse{
			ID:          bookstock.Book.ID,
			Title:       bookstock.Book.Title,
			Description: bookstock.Book.Description,
			CreatedAt:   bookstock.Book.CreatedAt,
			UpdatedAt:   bookstock.Book.UpdatedAt,
		}

		if bookstock.Book.Cover != nil {
			bookResponse.Cover = &dto.MediaResponse{
				ID:        bookstock.Book.Cover.ID,
				Path:      bookstock.Book.Cover.Path,
				CreatedAt: bookstock.Book.Cover.CreatedAt,
			}
		}

		response.Book = bookResponse
	}

	return response
}
