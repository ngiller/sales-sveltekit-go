package utils

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// SuccessResponse returns a standard JSON payload for successful requests
func SuccessResponse(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"data": data,
	})
}

// ErrorResponse returns a standard JSON payload for failed requests
func ErrorResponse(c *fiber.Ctx, status int, message interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
	})
}

// allowedExtensions maps MIME categories to valid extensions
var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}
var allowedDocumentExts = map[string]bool{
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
}

// ValidateFile checks file extension and size. Returns error message or empty string.
func ValidateFile(file *multipart.FileHeader, allowedExts map[string]bool, maxSize int64) string {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		return "File type not allowed"
	}
	if file.Size > maxSize {
		return "File size exceeds maximum limit"
	}
	return ""
}

// AllowedImageExts returns the map of allowed image extensions
func AllowedImageExts() map[string]bool {
	return allowedImageExts
}

// AllowedDocumentExts returns the map of allowed document extensions
func AllowedDocumentExts() map[string]bool {
	return allowedDocumentExts
}

// MergeExts combines multiple extension maps
func MergeExts(maps ...map[string]bool) map[string]bool {
	result := make(map[string]bool)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
