package middleware

import (
	"backend/internal/models"
	"backend/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RequirePolicy creates a middleware that checks if the logged-in user has the required policy action on the specified table
func RequirePolicy(db *gorm.DB, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(uint)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		var user models.User
		if err := db.Select("user_group_id").First(&user, userID).Error; err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
		}

		if user.UserGroupID == nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access Denied: User does not belong to any group")
		}

		// Admin bypass - group_id 1 (super admin) has all permissions
		if *user.UserGroupID == 1 {
			return c.Next()
		}

		// Extract endpoint from path (e.g. /api/payment-terms/1 -> payment-terms)
		segments := strings.Split(strings.Trim(c.Path(), "/"), "/")
		if len(segments) < 2 {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Invalid API path format")
		}
		endpoint := segments[1]

		var table models.MasterTableAccess
		if err := db.Where("endpoint = ?", endpoint).First(&table).Error; err != nil {
			// Jika endpoint tidak ditemukan di menu_access, tolak akses (default deny)
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Endpoint not registered for policy check")
		}

		var policyCount int64
		err := db.Model(&models.GroupPolicy{}).
			Where("group_id = ? AND action = ? AND (table_name = ? OR table_id = ?)", *user.UserGroupID, action, table.Name, table.ID).
			Count(&policyCount).Error

		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error during policy check")
		}

		if policyCount == 0 {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access Denied: You don't have permission to perform this action")
		}

		return c.Next()
	}
}
