package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DepartementHandler struct {
	repo *repository.DepartementRepository
}

func NewDepartementHandler(repo *repository.DepartementRepository) *DepartementHandler {
	return &DepartementHandler{repo: repo}
}

func (h *DepartementHandler) FindAll(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	sortBy := c.Query("sort")
	sortDir := c.Query("order")

	depts, total, err := h.repo.FindAll(search, page, limit, sortBy, sortDir)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve departements")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"departments": depts,
		"total":       total,
		"page":        page,
		"limit":       limit,
	})
}

func (h *DepartementHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	dept, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Departement not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, dept)
}

func (h *DepartementHandler) Create(c *fiber.Ctx) error {
	var dept models.Departement
	if err := c.BodyParser(&dept); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Set UserCreated from logged in user
	userID := c.Locals("user_id").(uint)
	dept.UserCreated = &userID
	dept.UserUpdate = &userID

	// Use transaction
	err := h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dept).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create departement: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, dept)
}

func (h *DepartementHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var reqDept models.Departement
	if err := c.BodyParser(&reqDept); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	dept, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Departement not found")
	}

	dept.Name = reqDept.Name
	dept.PropertyID = reqDept.PropertyID
	
	// Set UserUpdate from logged in user
	userID := c.Locals("user_id").(uint)
	dept.UserUpdate = &userID

	// Use transaction
	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(dept).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update departement: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, dept)
}

func (h *DepartementHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	// Use transaction
	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Departement{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete departement: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Departement successfuly deleted"})
}
