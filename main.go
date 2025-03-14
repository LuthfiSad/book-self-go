package main

import (
	"go-rest-api/internal/api"
	"go-rest-api/internal/config"
	"go-rest-api/internal/connection"
	"go-rest-api/internal/middleware"
	"go-rest-api/internal/repository"
	"go-rest-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()

	dbConnection, dbGorm := connection.GetDatabase(cnf.Database)

	userRepository := repository.NewUser(dbConnection)
	mediaRepository := repository.NewMediaRepositoryImpl(dbGorm)
	mediaService := service.NewMediaService(mediaRepository, cnf)
	bookRepository := repository.NewBookRepository(dbGorm)
	bookService := service.NewBookService(bookRepository, mediaRepository, cnf)

	authService := service.NewAuth(cnf, userRepository)

	authHandler := middleware.Authenticate(authService)

	app := fiber.New()

	api.NewAuth(app, authHandler, authService)
	api.NewBookApi(app, authHandler, bookService)
	api.NewMediaApi(app, authHandler, mediaService, cnf)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
