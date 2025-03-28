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
	bookRepository := repository.NewBookRepository(dbGorm)
	BookstockRepository := repository.NewBookstockRepositoryImpl(dbGorm)
	BookTransactionRepository := repository.NewBookTransactionRepositoryImpl(dbGorm)
	CustomerRepository := repository.NewCustomerRepositoryImpl(dbGorm)

	bookService := service.NewBookService(bookRepository, mediaRepository, cnf)
	mediaService := service.NewMediaService(mediaRepository, bookService, cnf)
	bookstockService := service.NewBookstockService(BookstockRepository, bookRepository)
	bookTransactionService := service.NewBookTransactionService(BookTransactionRepository, bookRepository, BookstockRepository, CustomerRepository)
	customerService := service.NewCustomerService(CustomerRepository)

	authService := service.NewAuth(cnf, userRepository)

	authHandler := middleware.Authenticate(authService)
	fileHandler := middleware.FileUploadMiddleware(cnf)

	app := fiber.New()

	api.NewAuth(app, authHandler, authService)
	api.NewBookApi(app, authHandler, bookService)
	api.NewMediaApi(app, authHandler, fileHandler, mediaService, cnf)
	api.NewBookstockApi(app, authHandler, bookstockService)
	api.NewBookTransactionApi(app, authHandler, bookTransactionService)
	api.NewCustomerApi(app, authHandler, customerService)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
