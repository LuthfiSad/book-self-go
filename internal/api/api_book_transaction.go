package api

import (
	"context"
	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type bookTransactionApi struct {
	bookTransactionService domain.BookTransactionService
}

func NewBookTransactionApi(app *fiber.App, authHandler fiber.Handler, bookTransactionService domain.BookTransactionService) {
	bta := bookTransactionApi{
		bookTransactionService: bookTransactionService,
	}

	bookTransactionGroup := app.Group("/v1/book-transactions")

	bookTransactionGroup.Get("/", authHandler, bta.getAllBookTransactions)
	bookTransactionGroup.Post("/", authHandler, bta.createBookTransaction)
	bookTransactionGroup.Put("/:id", authHandler, bta.updateBookTransaction)
	bookTransactionGroup.Put("/:id/return", authHandler, bta.returnBookTransaction)
	bookTransactionGroup.Delete("/:id", authHandler, bta.deleteBookTransaction)
}

func (bta *bookTransactionApi) getAllBookTransactions(ctx *fiber.Ctx) error {
	filters := make(map[string]interface{})

	if search := ctx.Query("search"); search != "" {
		filters["search"] = search
	}
	if status := ctx.Query("status"); status != "" {
		filters["status"] = status
	}

	transactions, err := bta.bookTransactionService.GetAllBookTransactions(filters)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(transactions))
}

func (bta *bookTransactionApi) createBookTransaction(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.BookTransactionCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(err.Error()))
	}

	validationErrors := utils.Validate(req)
	if len(validationErrors) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(validationErrors))
	}

	transaction, err := bta.bookTransactionService.CreateBookTransaction(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusCreated).JSON(dto.NewResponseData(transaction))
}

func (bta *bookTransactionApi) updateBookTransaction(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid ID format"))
	}

	var req dto.BookTransactionUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(err.Error()))
	}

	transaction, err := bta.bookTransactionService.UpdateBookTransaction(c, id, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(transaction))
}

func (bta *bookTransactionApi) returnBookTransaction(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.BookTransactionUpdateStatusRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid request body"))
	}

	transaction, err := bta.bookTransactionService.ReturnBookTransaction(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(transaction))
}

func (bta *bookTransactionApi) deleteBookTransaction(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid ID format"))
	}

	if err := bta.bookTransactionService.DeleteBookTransaction(c, id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseMessage("Book transaction deleted successfully"))
}
