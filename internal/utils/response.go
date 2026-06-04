package utils

import "github.com/gofiber/fiber/v2"

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
