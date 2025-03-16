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

type bookApi struct {
	bookService domain.BookService
}

func NewBookApi(app *fiber.App, authHandler fiber.Handler, bookService domain.BookService) {
	ba := bookApi{
		bookService: bookService,
	}

	bookGroup := app.Group("/v1/books")

	bookGroup.Get("/", authHandler, ba.getAllBooks)
	bookGroup.Get("/:id", authHandler, ba.getBookByID)
	bookGroup.Post("/", authHandler, ba.createBook)
	bookGroup.Put("/:id", authHandler, ba.updateBook)
	bookGroup.Delete("/:id", authHandler, ba.deleteBook)
	bookGroup.Delete("/cover/:id", authHandler, ba.deleteBookCover)
}

func (ba *bookApi) getAllBooks(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	books, err := ba.bookService.GetAllBooks(c)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(books))
}

func (ba *bookApi) getBookByID(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid ID format"))
	}

	book, err := ba.bookService.GetBookByID(c, id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dto.NewResponseMessage("Book not found"))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(book))
}

func (ba *bookApi) getBookByCoverID(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid ID format"))
	}

	book, err := ba.bookService.GetBookByCoverID(c, id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(book))
}

func (ba *bookApi) createBook(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.BookCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid request body"))
	}

	validationErrors := utils.Validate(req)
	if len(validationErrors) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(validationErrors))
	}

	book, err := ba.bookService.CreateBook(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage("Failed to create book"))
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.NewResponseData(book))
}

func (ba *bookApi) updateBook(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid ID format"))
	}

	var req dto.BookUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(err.Error()))
	}

	book, err := ba.bookService.UpdateBook(c, id, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.NewResponseData(book))
}

func (ba *bookApi) deleteBook(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid ID format"))
	}

	if err := ba.bookService.DeleteBook(c, id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage("Failed to delete book"))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseMessage("Book deleted successfully"))
}

func (ba *bookApi) deleteBookCover(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid ID format"))
	}

	if err := ba.bookService.DeleteBookCover(c, id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseMessage("Book cover deleted successfully"))
}
