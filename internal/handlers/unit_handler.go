package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UnitHandler struct {
	repo *repository.UnitRepository
}

func NewUnitHandler(repo *repository.UnitRepository) *UnitHandler {
	return &UnitHandler{repo: repo}
}

func (h *UnitHandler) FindAll(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	sortBy := c.Query("sort")
	sortDir := c.Query("order")

	items, total, err := h.repo.FindAll(search, page, limit, sortBy, sortDir)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve units")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *UnitHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	item, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Unit not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, item)
}

func (h *UnitHandler) Create(c *fiber.Ctx) error {
	var item models.Unit
	if err := c.BodyParser(&item); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	userID := c.Locals("user_id").(uint)
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	item.UserCreated = &userIDStr
	item.UserUpdate = &userIDStr

	err := h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create unit: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, item)
}

func (h *UnitHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var reqItem models.Unit
	if err := c.BodyParser(&reqItem); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	item, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Unit not found")
	}

	item.Name = reqItem.Name
	if reqItem.PropertyID != nil {
		item.PropertyID = reqItem.PropertyID
	}

	userID := c.Locals("user_id").(uint)
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	item.UserUpdate = &userIDStr

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(item).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update unit: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, item)
}

func (h *UnitHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Unit{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete unit: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Unit successfully deleted"})
}
