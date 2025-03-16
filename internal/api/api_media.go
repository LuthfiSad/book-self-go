package api

import (
	"context"
	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/config"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mediaApi struct {
	mediaService domain.MediaService
	cnf          *config.Config
}

func NewMediaApi(app *fiber.App, authHandler fiber.Handler, fileHandler fiber.Handler, mediaService domain.MediaService, cnf *config.Config) {
	ma := mediaApi{
		mediaService: mediaService,
		cnf:          cnf,
	}

	mediaGroup := app.Group("/v1/media")

	mediaGroup.Get("/", authHandler, ma.getAllMedia)
	mediaGroup.Get("/:id", authHandler, ma.getMediaByID)
	mediaGroup.Get("/file/:name", authHandler, ma.getMediaByName)
	mediaGroup.Post("/", authHandler, fileHandler, ma.uploadMedia)
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
	fileName := c.Locals("fileName")
	if fileName == nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewResponseMessage("No file provided"))
	}

	filePath := c.Locals("filePath")
	if filePath == nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewResponseMessage("No file path provided"))
	}

	fileNameStr := fileName.(string)
	filePathStr := filePath.(string)

	mediaResponse, err := ma.mediaService.UploadMedia(fileNameStr, filePathStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.NewResponseData(mediaResponse))
}

func (ma *mediaApi) deleteMedia(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid ID format"))
	}

	if err := ma.mediaService.DeleteMedia(c, id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.NewResponseMessage(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewResponseMessage("Media deleted successfully"))
}
