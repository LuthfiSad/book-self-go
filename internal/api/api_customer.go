package api

import (
	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CustomerApi struct {
	customerService domain.CustomerService
}

func NewCustomerApi(app *fiber.App, authHandler fiber.Handler, customerService domain.CustomerService) {
	ch := CustomerApi{customerService: customerService}

	CustomerGroup := app.Group("/v1/customers")

	CustomerGroup.Get("/", authHandler, ch.GetAllCustomers)
	CustomerGroup.Get("/:id", authHandler, ch.GetCustomerByID)
	CustomerGroup.Post("/", authHandler, ch.CreateCustomer)
	CustomerGroup.Put("/:id", authHandler, ch.UpdateCustomer)
	CustomerGroup.Delete("/:id", authHandler, ch.DeleteCustomer)
}

func (h *CustomerApi) GetAllCustomers(ctx *fiber.Ctx) error {
	customers, err := h.customerService.GetAllCustomers()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}
	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(customers))
}

func (h *CustomerApi) GetCustomerByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid customer ID"))
	}

	customer, err := h.customerService.GetCustomerByID(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dto.NewResponseMessage("Customer not found"))
	}
	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(customer))
}

func (h *CustomerApi) CreateCustomer(ctx *fiber.Ctx) error {
	var req dto.CustomerCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid request body"))
	}

	validationErrors := utils.Validate(req)
	if len(validationErrors) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(validationErrors))
	}

	customer, err := h.customerService.CreateCustomer(req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusCreated).JSON(dto.NewResponseData(customer))
}

func (h *CustomerApi) UpdateCustomer(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid customer ID"))
	}

	var req dto.CustomerUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid request body"))
	}

	customer, err := h.customerService.UpdateCustomer(id, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData(customer))
}

func (h *CustomerApi) DeleteCustomer(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid customer ID"))
	}

	err = h.customerService.DeleteCustomer(id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseMessage("Customer deleted successfully"))
}
