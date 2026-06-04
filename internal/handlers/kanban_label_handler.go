package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type KanbanLabelHandler struct {
	repo *repository.KanbanLabelRepository
}

func NewKanbanLabelHandler(repo *repository.KanbanLabelRepository) *KanbanLabelHandler {
	return &KanbanLabelHandler{repo: repo}
}

func (h *KanbanLabelHandler) FindByBoardID(c *fiber.Ctx) error {
	boardIDStr := c.Query("board_id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "board_id is required")
	}

	labels, err := h.repo.FindByBoardID(uint(boardID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve labels")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, labels)
}

func (h *KanbanLabelHandler) Create(c *fiber.Ctx) error {
	var req struct {
		BoardID uint   `json:"board_id"`
		Name    string `json:"name"`
		Color   string `json:"color"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" || req.Color == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Name and color are required")
	}

	label := &models.KanbanLabel{
		BoardID: req.BoardID,
		Name:    req.Name,
		Color:   req.Color,
	}

	err := h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(label).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create label: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, label)
}

func (h *KanbanLabelHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	label, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Label not found")
	}

	if req.Name != "" {
		label.Name = req.Name
	}
	if req.Color != "" {
		label.Color = req.Color
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Save(label).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update label: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, label)
}

func (h *KanbanLabelHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&models.KanbanLabel{}, id).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete label: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Label successfully deleted"})
}
