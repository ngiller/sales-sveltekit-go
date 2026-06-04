package handlers

import (
	"backend/internal/repository"
	"backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type ProductCategoryHandler struct {
	repo *repository.ProductCategoryRepository
}

func NewProductCategoryHandler(repo *repository.ProductCategoryRepository) *ProductCategoryHandler {
	return &ProductCategoryHandler{repo: repo}
}

func (h *ProductCategoryHandler) FindAll(c *fiber.Ctx) error {
	categories, err := h.repo.FindAll()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve product categories: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"categories": categories,
	})
}
