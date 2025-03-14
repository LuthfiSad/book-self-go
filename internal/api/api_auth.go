package api

import (
	"context"
	"errors"
	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/constants"
	"go-rest-api/internal/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type authApi struct {
	authService domain.AuthService
}

func NewAuth(app *fiber.App, authHandler fiber.Handler,
	authService domain.AuthService) {

	ha := authApi{
		authService: authService,
	}

	authenticateGroup := app.Group("/v1/authenticate")

	authenticateGroup.Post("/", ha.authenticate)
	authenticateGroup.Post("/register", ha.register)
	authenticateGroup.Post("/validate", authHandler, ha.authenticateValidate)
}

func (a authApi) register(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.RegisterReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.NewResponseMessage("Invalid request format"))
	}

	validationErrors := utils.Validate(req)
	if len(validationErrors) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(validationErrors))
	}

	res, err := a.authService.Register(c, req)
	if err != nil {
		if errors.Is(err, constants.ErrEmailAlreadyExists) {
			return ctx.Status(http.StatusConflict).JSON(dto.NewResponseMessage("Email is already registered"))
		}
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage("Internal server error"))
	}

	return ctx.Status(http.StatusCreated).JSON(dto.NewResponseData[dto.UserData](res))
}

func (a authApi) authenticate(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.AuthReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	validationErrors := utils.Validate(req)
	if len(validationErrors) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(validationErrors))
	}

	res, err := a.authService.Authenticate(c, req)

	if err != nil {
		if errors.Is(err, constants.ErrInvalidCredential) {
			return ctx.Status(http.StatusUnauthorized).JSON(dto.NewResponseMessage("Invalid credential. Please check your username and password, then try again. If the problem persists, contact support for assistance."))
		}
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage("An internal server error has occurred. Please try again later. If the issue persists, contact support for further assistance. We apologize for any inconvenience."))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData[dto.AuthRes](res))
}

func (a authApi) authenticateValidate(ctx *fiber.Ctx) error {
	userLocal := ctx.Locals("x-user")
	if userLocal == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage("Sorry, the token you entered is invalid. Please check your token and try again or contact customer support for further assistance. Thank you."))
	}
	return ctx.Status(http.StatusOK).JSON(dto.NewResponseData[dto.UserData](userLocal.(dto.UserData)))
}
