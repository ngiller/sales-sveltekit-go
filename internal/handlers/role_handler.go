package handlers

import (
	"encoding/json"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RoleHandler struct {
	repo *repository.RoleRepository
}

func NewRoleHandler(repo *repository.RoleRepository) *RoleHandler {
	return &RoleHandler{repo: repo}
}

func (h *RoleHandler) checkViewAllPolicy(groupID uint) bool {
	var table models.MasterTableAccess
	if err := h.repo.GetDB().Where("endpoint = ?", "quotations").First(&table).Error; err != nil {
		return false
	}
	var count int64
	h.repo.GetDB().Model(&models.GroupPolicy{}).
		Where("group_id = ? AND action = ? AND (table_name = ? OR table_id = ?)",
			groupID, "view_all", table.Name, table.ID).
		Count(&count)
	return count > 0
}

func (h *RoleHandler) FindAll(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	sortBy := c.Query("sort")
	sortDir := c.Query("order")

	roles, total, err := h.repo.FindAll(search, page, limit, sortBy, sortDir)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve roles")
	}

	type roleResponse struct {
		models.UserGroup
		CanViewAllQuotations bool `json:"can_view_all_quotations"`
	}
	items := make([]roleResponse, len(roles))
	for i, r := range roles {
		items[i] = roleResponse{
			UserGroup:            r,
			CanViewAllQuotations: h.checkViewAllPolicy(r.ID),
		}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"roles": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *RoleHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	role, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Role not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, role)
}

func (h *RoleHandler) setViewAllPolicy(groupID uint, enabled bool) error {
	var table models.MasterTableAccess
	if err := h.repo.GetDB().Where("endpoint = ?", "quotations").First(&table).Error; err != nil {
		return err
	}

	if enabled {
		var count int64
		h.repo.GetDB().Model(&models.GroupPolicy{}).
			Where("group_id = ? AND action = ? AND (table_name = ? OR table_id = ?)",
				groupID, "view_all", table.Name, table.ID).
			Count(&count)
		if count == 0 {
			policy := models.GroupPolicy{
				GroupID:         groupID,
				TargetTableName: table.Name,
				TableID:         table.ID,
				Action:          "view_all",
			}
			return h.repo.GetDB().Create(&policy).Error
		}
	} else {
		return h.repo.GetDB().Where("group_id = ? AND action = ? AND (table_name = ? OR table_id = ?)",
			groupID, "view_all", table.Name, table.ID).
			Delete(&models.GroupPolicy{}).Error
	}
	return nil
}

func (h *RoleHandler) Create(c *fiber.Ctx) error {
	var role models.UserGroup
	if err := c.BodyParser(&role); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Extract can_view_all_quotations from raw body
	var bodyMap map[string]interface{}
	json.Unmarshal(c.Body(), &bodyMap)
	canViewAll, _ := bodyMap["can_view_all_quotations"].(bool)

	// Set audit fields
	userID := c.Locals("user_id").(uint)
	role.UserCreated = &userID
	role.UserUpdate = &userID

	// Use transaction
	err := h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create role: "+err.Error())
	}

	// Manage view_all policy
	if err := h.setViewAllPolicy(role.ID, canViewAll); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to set policy: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, role)
}

func (h *RoleHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var reqRole models.UserGroup
	if err := c.BodyParser(&reqRole); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Extract can_view_all_quotations from raw body
	var bodyMap map[string]interface{}
	json.Unmarshal(c.Body(), &bodyMap)
	canViewAll, _ := bodyMap["can_view_all_quotations"].(bool)

	role, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Role not found")
	}

	role.Name = reqRole.Name

	// Update audit fields
	userID := c.Locals("user_id").(uint)
	role.UserUpdate = &userID

	// Use transaction
	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(role).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update role: "+err.Error())
	}

	// Manage view_all policy
	if err := h.setViewAllPolicy(uint(id), canViewAll); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to set policy: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, role)
}

func (h *RoleHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	// Use transaction
	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.UserGroup{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete role: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Role successfully deleted"})
}
