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

	dbConnection := connection.GetDatabase(cnf.Database)

	userRepository := repository.NewUser(dbConnection)
	authService := service.NewAuth(cnf, userRepository)

	authHandler := middleware.Authenticate(authService)

	app := fiber.New()

	api.NewAuth(app, authHandler, authService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
