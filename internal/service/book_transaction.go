package service

import (
	"context"
	"errors"
	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/constants"
	"go-rest-api/internal/repository"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type bookTransactionService struct {
	bookTransactionRepo domain.BookTransactionRepository
	bookRepo            domain.BookRepository
	bookstockRepo       domain.BookstockRepository
	customerRepo        domain.CustomerRepository
}

func NewBookTransactionService(
	bookTransactionRepo domain.BookTransactionRepository,
	bookRepo domain.BookRepository,
	bookstockRepo domain.BookstockRepository,
	customerRepo domain.CustomerRepository,
) domain.BookTransactionService {
	return &bookTransactionService{
		bookTransactionRepo: bookTransactionRepo,
		bookRepo:            bookRepo,
		bookstockRepo:       bookstockRepo,
		customerRepo:        customerRepo,
	}
}

func (s *bookTransactionService) GetAllBookTransactions(filter map[string]interface{}) ([]dto.BookTransactionResponse, error) {
	book_transactions, err := s.bookTransactionRepo.Find(filter)
	if err != nil {
		return nil, err
	}

	bookTransactionResponses := make([]dto.BookTransactionResponse, 0, len(book_transactions))
	for _, book_transaction := range book_transactions {
		bookTransactionResponses = append(bookTransactionResponses, s.toBookTransactionResponse(&book_transaction))
	}

	return bookTransactionResponses, nil
}

func (s *bookTransactionService) CreateBookTransaction(ctx context.Context, req dto.BookTransactionCreateRequest) (*dto.BookTransactionResponse, error) {
	book, err := s.bookRepo.FindByID(ctx, req.BookID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errors.New("invalid book ID: book not found")
	}

	bookstock, err := s.bookstockRepo.FindByCode(req.StockCode)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errors.New("invalid stock code: stock not found")
	}

	if bookstock.Status != constants.BookStockStatusAvailable {
		return nil, errors.New("book stock is not available for borrowing")
	}

	customer, err := s.customerRepo.FindByID(req.CustomerID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errors.New("invalid customer ID: customer not found")
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errors.New("invalid due date format: use YYYY-MM-DD")
	}

	now := time.Now()
	book_transaction := &domain.BookTransaction{
		ID:         uuid.New(),
		BookID:     req.BookID,
		Book:       *book,
		StockCode:  req.StockCode,
		BookStock:  *bookstock,
		CustomerID: req.CustomerID,
		Customer:   *customer,
		DueDate:    dueDate,
		Status:     constants.BookTransactionStatusBorrowed,
		BorrowedAt: &now,
		ReturnAt:   nil,
	}

	// Begin transaction
	tx := s.bookTransactionRepo.(*repository.BookTransactionRepositoryImpl).GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(book_transaction).Error; err != nil {
		tx.Rollback()
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	bookstock.Status = constants.BookTransactionStatusBorrowed
	bookstock.BorrowedID = &req.CustomerID
	bookstock.BorrowedAt = &now

	if err := tx.Save(bookstock).Error; err != nil {
		tx.Rollback()
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	response := s.toBookTransactionResponse(book_transaction)
	return &response, nil
}

func (s *bookTransactionService) UpdateBookTransaction(ctx context.Context, id uuid.UUID, req dto.BookTransactionUpdateRequest) (*dto.BookTransactionResponse, error) {
	bookTransaction, err := s.bookTransactionRepo.FindByID(id)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errors.New("book transaction not found")
	}

	if req.BookID != uuid.Nil {
		bookTransaction.BookID = req.BookID
	}
	if req.StockCode != "" {
		bookTransaction.StockCode = req.StockCode
	}
	if req.CustomerID != uuid.Nil {
		bookTransaction.CustomerID = req.CustomerID
	}
	if req.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			return nil, errors.New("invalid due date format, use YYYY-MM-DD")
		}
		bookTransaction.DueDate = parsedDate
	}
	if req.Status != "" {
		bookTransaction.Status = req.Status
	}

	if err := s.bookTransactionRepo.Update(bookTransaction); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	response := s.toBookTransactionResponse(bookTransaction)
	return &response, nil
}

func (s *bookTransactionService) ReturnBookTransaction(ctx context.Context, req dto.BookTransactionUpdateStatusRequest) (*dto.BookTransactionResponse, error) {
	// Find book_transaction
	book_transaction, err := s.bookTransactionRepo.FindByID(req.ID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errors.New("book_transaction not found")
	}

	if book_transaction.Status != constants.BookTransactionStatusBorrowed {
		return nil, errors.New("book_transaction is not in borrowed status")
	}

	// Find bookstock
	bookstock, err := s.bookstockRepo.FindByCode(book_transaction.StockCode)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errors.New("book stock not found")
	}

	now := time.Now()

	// Begin transaction
	tx := s.bookTransactionRepo.(*repository.BookTransactionRepositoryImpl).GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	book_transaction.Status = constants.BookTransactionStatusAvailable
	book_transaction.ReturnAt = &now

	if err := tx.Save(book_transaction).Error; err != nil {
		tx.Rollback()
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	bookstock.Status = constants.BookStockStatusAvailable
	bookstock.BorrowedID = nil
	bookstock.BorrowedAt = nil

	if err := tx.Save(bookstock).Error; err != nil {
		tx.Rollback()
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	response := s.toBookTransactionResponse(book_transaction)
	return &response, nil
}

func (s *bookTransactionService) DeleteBookTransaction(ctx context.Context, id uuid.UUID) error {
	_, err := s.bookTransactionRepo.FindByID(id)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return errors.New("book_transaction not found")
	}

	if err := s.bookTransactionRepo.Delete(id); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	return nil
}

func (s *bookTransactionService) toBookTransactionResponse(book_transaction *domain.BookTransaction) dto.BookTransactionResponse {
	response := dto.BookTransactionResponse{
		ID:         book_transaction.ID,
		BookID:     book_transaction.BookID,
		StockCode:  book_transaction.StockCode,
		CustomerID: book_transaction.CustomerID,
		DueDate:    book_transaction.DueDate,
		Status:     book_transaction.Status,
		BorrowedAt: book_transaction.BorrowedAt,
		ReturnAt:   book_transaction.ReturnAt,
	}

	bookResponse := &dto.BookResponse{
		ID:          book_transaction.Book.ID,
		Title:       book_transaction.Book.Title,
		Description: book_transaction.Book.Description,
		CreatedAt:   book_transaction.Book.CreatedAt,
		UpdatedAt:   book_transaction.Book.UpdatedAt,
	}

	if book_transaction.Book.Cover != nil {
		bookResponse.Cover = &dto.MediaResponse{
			ID:        book_transaction.Book.Cover.ID,
			Path:      book_transaction.Book.Cover.Path,
			CreatedAt: book_transaction.Book.Cover.CreatedAt,
		}
	}

	response.Book = bookResponse

	response.BookStock = &dto.BookstockResponse{
		Code:       book_transaction.BookStock.Code,
		BookID:     book_transaction.BookStock.BookID,
		Status:     book_transaction.BookStock.Status,
		BorrowedID: book_transaction.BookStock.BorrowedID,
		BorrowedAt: book_transaction.BookStock.BorrowedAt,
	}

	response.Customer = &dto.CustomerResponse{
		ID:        book_transaction.Customer.ID,
		Name:      book_transaction.Customer.Name,
		CreatedAt: book_transaction.Customer.CreatedAt,
		UpdatedAt: book_transaction.Customer.UpdatedAt,
	}

	return response
}
