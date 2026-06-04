package handlers

import (
	"backend/internal/repository"
	"backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type BrandHandler struct {
	repo *repository.BrandRepository
}

func NewBrandHandler(repo *repository.BrandRepository) *BrandHandler {
	return &BrandHandler{repo: repo}
}

func (h *BrandHandler) FindAll(c *fiber.Ctx) error {
	brands, err := h.repo.FindAll()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve brands: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"brands": brands,
	})
}
