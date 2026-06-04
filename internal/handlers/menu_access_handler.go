package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MenuAccessHandler struct {
	repo     *repository.MenuAccessRepository
	menuRepo *repository.MenuRepository
	db       *gorm.DB
}

func NewMenuAccessHandler(repo *repository.MenuAccessRepository, menuRepo *repository.MenuRepository, db *gorm.DB) *MenuAccessHandler {
	return &MenuAccessHandler{repo: repo, menuRepo: menuRepo, db: db}
}

func (h *MenuAccessHandler) FindAll(c *fiber.Ctx) error {
	items, err := h.repo.FindAll()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve menu navigation")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, items)
}

func (h *MenuAccessHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	item, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Menu navigation not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, item)
}

func (h *MenuAccessHandler) Create(c *fiber.Ctx) error {
	var item models.MasterTableAccess
	if err := c.BodyParser(&item); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.repo.Create(&item); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create menu navigation")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, item)
}

func (h *MenuAccessHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var reqItem models.MasterTableAccess
	if err := c.BodyParser(&reqItem); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	item, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Menu navigation not found")
	}

	item.Name = reqItem.Name
	item.ParentID = reqItem.ParentID
	item.MenuName = reqItem.MenuName
	item.Path = reqItem.Path
	item.Endpoint = reqItem.Endpoint
	item.Icon = reqItem.Icon
	item.SortOrder = reqItem.SortOrder
	item.IsActive = reqItem.IsActive

	if err := h.repo.Update(item); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update menu navigation")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, item)
}

func (h *MenuAccessHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete menu navigation")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Menu navigation successfully deleted"})
}
