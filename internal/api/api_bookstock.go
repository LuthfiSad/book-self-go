package api

import (
	"context"
	"net/http"
	"time"

	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type bookstockApi struct {
	bookstockService domain.BookstockService
}

func NewBookstockApi(app *fiber.App, authHandler fiber.Handler, bookstockService domain.BookstockService) {
	ba := bookstockApi{
		bookstockService: bookstockService,
	}

	bookstockGroup := app.Group("/v1/bookstocks")

	bookstockGroup.Get("/", authHandler, ba.getAllBookstocks)
	bookstockGroup.Get("/:code", authHandler, ba.getBookstockByCode)
	bookstockGroup.Get("/book/:bookId", authHandler, ba.getBookstocksByBookID)
	bookstockGroup.Get("/book/:bookId/available", authHandler, ba.getAvailableBookstocksByBookID)
	bookstockGroup.Post("/", authHandler, ba.createBookstock)
	bookstockGroup.Put("/:code", authHandler, ba.updateBookstock)
	bookstockGroup.Delete("/:code", authHandler, ba.deleteBookstock)
}

func (ba *bookstockApi) getAllBookstocks(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	_ = c // Using the timeout context even though the service method doesn't take it yet

	bookstocks, err := ba.bookstockService.GetAllBookstocks()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(bookstocks))
}

func (ba *bookstockApi) getBookstockByCode(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	_ = c // Using the timeout context

	code := ctx.Params("code")
	if code == "" {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid code parameter"))
	}

	bookstock, err := ba.bookstockService.GetBookstockByCode(code)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dto.NewResponseMessage("Bookstock not found"))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(bookstock))
}

func (ba *bookstockApi) getBookstocksByBookID(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	_ = c // Using the timeout context

	bookId, err := uuid.Parse(ctx.Params("bookId"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid book ID format"))
	}

	bookstocks, err := ba.bookstockService.GetBookstocksByBookID(bookId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(bookstocks))
}

func (ba *bookstockApi) getAvailableBookstocksByBookID(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	_ = c // Using the timeout context

	bookId, err := uuid.Parse(ctx.Params("bookId"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid book ID format"))
	}

	bookstocks, err := ba.bookstockService.GetAvailableBookstocksByBookID(bookId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(bookstocks))
}

func (ba *bookstockApi) createBookstock(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	_ = c // Using the timeout context

	var req dto.BookstockCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(err.Error()))
	}

	validationErrors := utils.Validate(req)
	if len(validationErrors) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(validationErrors))
	}

	bookstock, err := ba.bookstockService.CreateBookstock(req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.NewResponseData(bookstock))
}

func (ba *bookstockApi) updateBookstock(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	_ = c // Using the timeout context

	code := ctx.Params("code")
	if code == "" {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid code parameter"))
	}

	var req dto.BookstockUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(err.Error()))
	}

	validationErrors := utils.Validate(req)
	if len(validationErrors) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(validationErrors))
	}

	bookstock, err := ba.bookstockService.UpdateBookstock(code, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.NewResponseData(bookstock))
}

func (ba *bookstockApi) deleteBookstock(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	_ = c // Using the timeout context

	code := ctx.Params("code")
	if code == "" {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid code parameter"))
	}

	if err := ba.bookstockService.DeleteBookstock(code); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseMessage("Bookstock deleted successfully"))
}
