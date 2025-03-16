package middleware

import (
	"fmt"
	"go-rest-api/dto"
	"go-rest-api/internal/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func FileUploadMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request is multipart/form-data
		contentType := c.Get("Content-Type")
		if !strings.Contains(contentType, "multipart/form-data") {
			return c.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Content-Type must be multipart/form-data"))
		}

		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("File upload is required"))
		}

		maxUploadSize, err := strconv.ParseInt(cfg.File.MaxUploadSize, 10, 64)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid max upload size"))
		}

		if file.Size > maxUploadSize*1024*1024 {
			return c.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage(fmt.Sprintf("File size exceeds limit of %d MB", maxUploadSize)))
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
		}

		if !allowedExts[strings.ToLower(ext)] {
			return c.Status(http.StatusBadRequest).JSON(dto.NewResponseMessage("Invalid file type. Allowed: jpg, jpeg, png, gif"))
		}

		uploadPath := cfg.File.UploadPath
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(fmt.Errorf("failed to create upload directory: %w", err)))
		}

		newFileName := uuid.New().String() + ext
		filePath := filepath.Join(uploadPath, newFileName)

		src, err := file.Open()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(fmt.Errorf("failed to open uploaded file: %w", err)))
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(fmt.Errorf("failed to create destination file: %w", err)))
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(dto.NewResponseMessage(fmt.Errorf("failed to copy file: %w", err)))
		}

		c.Locals("fileName", cfg.File.LinkCover+"/"+newFileName)
		c.Locals("filePath", filePath)

		return c.Next()
	}
}
