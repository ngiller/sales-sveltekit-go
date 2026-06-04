package middleware

import (
	"backend/internal/models"
	"backend/internal/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RequireRole(db *gorm.DB, roleID uint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(uint)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		var user models.User
		if err := db.Select("user_group_id").First(&user, userID).Error; err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
		}

		if user.UserGroupID == nil || *user.UserGroupID != roleID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, fmt.Sprintf("Access Denied: Only users with role ID %d can access this resource", roleID))
		}

		return c.Next()
	}
}
