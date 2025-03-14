package middleware

import (
	"fmt"
	"go-rest-api/internal/config"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func FileUploadMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request is multipart/form-data
		if !strings.Contains(c.Get("Content-Type"), "multipart/form-data") {
			return c.Next()
		}

		// Get file from form
		file, err := c.FormFile("file")
		if err != nil {
			// If no file is provided, just continue
			if err == fiber.ErrUnprocessableEntity {
				return c.Next()
			}
			return err
		}

		maxUploadSize, err := strconv.ParseInt(cfg.File.MaxUploadSize, 10, 64)
		if err != nil {
			// handle error
		}

		if file.Size > maxUploadSize*1024*1024 {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("File size exceeds limit of %d MB", maxUploadSize))
		}

		// Check file extension
		ext := filepath.Ext(file.Filename)
		allowedExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
			".pdf":  true,
		}

		if !allowedExts[strings.ToLower(ext)] {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid file type")
		}

		// Add file to context locals
		c.Locals("file", file)

		return c.Next()
	}
}
