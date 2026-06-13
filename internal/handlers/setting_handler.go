package handlers

import (
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SettingHandler struct {
	repo *repository.SettingRepository
}

func NewSettingHandler(repo *repository.SettingRepository) *SettingHandler {
	return &SettingHandler{repo: repo}
}

func (h *SettingHandler) GetByCode(c *fiber.Ctx) error {
	codeStr := c.Params("code")
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid code")
	}

	propertyID := 1
	if pidStr := c.Query("property_id"); pidStr != "" {
		if pid, err := strconv.Atoi(pidStr); err == nil {
			propertyID = pid
		}
	}

	setting, err := h.repo.FindByCode(code, propertyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Setting not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, setting)
}
