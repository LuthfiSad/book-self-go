package api

import (
	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/config"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mediaApi struct {
	mediaService domain.MediaService
	cnf          *config.Config
}

func NewMediaApi(app *fiber.App, authHandler fiber.Handler, mediaService domain.MediaService, cnf *config.Config) {
	ma := mediaApi{
		mediaService: mediaService,
		cnf:          cnf,
	}

	mediaGroup := app.Group("/v1/media")

	mediaGroup.Get("/", authHandler, ma.getAllMedia)
	mediaGroup.Get("/:id", authHandler, ma.getMediaByID)
	mediaGroup.Get("/file/:name", authHandler, ma.getMediaByName)
	mediaGroup.Post("/", authHandler, ma.uploadMedia)
	mediaGroup.Delete("/:id", authHandler, ma.deleteMedia)
}

func (ma *mediaApi) getAllMedia(c *fiber.Ctx) error {
	media, err := ma.mediaService.GetAllMedia()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage("Failed to get media"))
	}

	return c.Status(http.StatusOK).JSON(dto.NewResponseData(media))
}

func (ma *mediaApi) getMediaByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid ID format"))
	}

	media, err := ma.mediaService.GetMediaByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(dto.NewResponseMessage("Media not found"))
	}

	return c.Status(http.StatusOK).JSON(dto.NewResponseData(media))
}

func (ma *mediaApi) getMediaByName(c *fiber.Ctx) error {
	fileName := c.Params("name")

	filePath := filepath.Join(ma.cnf.File.UploadPath, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(dto.NewResponseMessage("File not found"))
	}

	return c.SendFile(filePath)
}

func (ma *mediaApi) uploadMedia(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewResponseMessage("No file uploaded"))
	}

	mediaResponse, err := ma.mediaService.UploadMedia(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.NewResponseData(mediaResponse))
}

func (ma *mediaApi) deleteMedia(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid ID format"))
	}

	if err := ma.mediaService.DeleteMedia(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewResponseMessage("Failed to delete media"))
	}

	return c.Status(http.StatusOK).JSON(dto.NewResponseMessage("Media deleted successfully"))
}
