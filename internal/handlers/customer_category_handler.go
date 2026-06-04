package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CustomerCategoryHandler struct {
	repo *repository.CustomerCategoryRepository
}

func NewCustomerCategoryHandler(repo *repository.CustomerCategoryRepository) *CustomerCategoryHandler {
	return &CustomerCategoryHandler{repo: repo}
}

func (h *CustomerCategoryHandler) FindAll(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	sortBy := c.Query("sort")
	sortDir := c.Query("order")

	categories, total, err := h.repo.FindAll(search, page, limit, sortBy, sortDir)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve customer categories")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"categories": categories,
		"total":      total,
		"page":       page,
		"limit":      limit,
	})
}

func (h *CustomerCategoryHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	category, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Customer category not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, category)
}

func (h *CustomerCategoryHandler) Create(c *fiber.Ctx) error {
	var category models.CustomerCategory
	if err := c.BodyParser(&category); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Audit fields
	userID := c.Locals("user_id").(uint)
	category.UserCreated = &userID
	category.UserUpdate = &userID

	// Use transaction
	err := h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&category).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create customer category: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, category)
}

func (h *CustomerCategoryHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var reqCategory models.CustomerCategory
	if err := c.BodyParser(&reqCategory); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	category, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Customer category not found")
	}

	category.Name = reqCategory.Name
	
	// Audit fields
	userID := c.Locals("user_id").(uint)
	category.UserUpdate = &userID

	// Use transaction
	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(category).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update customer category: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, category)
}

func (h *CustomerCategoryHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	// Use transaction
	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.CustomerCategory{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete customer category: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Customer category successfully deleted"})
}
