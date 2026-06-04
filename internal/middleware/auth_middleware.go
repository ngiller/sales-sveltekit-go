package middleware

import (
	"backend/config"
	"backend/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			authHeader = c.Cookies("token")
		}

		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing or invalid JWT in request")
		}

		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		claims, err := config.ValidateJWT(tokenString)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		userID := uint(claims["user_id"].(float64))
		c.Locals("user_id", userID)
		c.Locals("email", claims["email"])
		c.Locals("inisial", claims["inisial"])

		return c.Next()
	}
}
