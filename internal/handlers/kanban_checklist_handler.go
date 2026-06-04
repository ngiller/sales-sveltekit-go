package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type KanbanChecklistHandler struct {
	repo *repository.KanbanChecklistRepository
}

func NewKanbanChecklistHandler(repo *repository.KanbanChecklistRepository) *KanbanChecklistHandler {
	return &KanbanChecklistHandler{repo: repo}
}

// Checklists

func (h *KanbanChecklistHandler) FindByCardID(c *fiber.Ctx) error {
	cardIDStr := c.Query("card_id")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "card_id is required")
	}

	checklists, err := h.repo.FindByCardID(uint(cardID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve checklists")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, checklists)
}

func (h *KanbanChecklistHandler) CreateChecklist(c *fiber.Ctx) error {
	var req struct {
		CardID uint   `json:"card_id"`
		Name   string `json:"name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Checklist name is required")
	}

	pos, err := h.repo.GetMaxChecklistPosition(req.CardID)
	if err != nil {
		pos = 0
	}

	cl := &models.KanbanChecklist{
		CardID:   req.CardID,
		Name:     req.Name,
		Position: pos,
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(cl).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create checklist: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, cl)
}

func (h *KanbanChecklistHandler) UpdateChecklist(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	cl, err := h.repo.FindChecklistByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Checklist not found")
	}

	if req.Name != "" {
		cl.Name = req.Name
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Save(cl).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update checklist: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, cl)
}

func (h *KanbanChecklistHandler) DeleteChecklist(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&models.KanbanChecklist{}, id).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete checklist: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Checklist deleted"})
}

// Checklist Items

func (h *KanbanChecklistHandler) CreateItem(c *fiber.Ctx) error {
	var req struct {
		ChecklistID uint   `json:"checklist_id"`
		Name        string `json:"name"`
		AssigneeID  *uint  `json:"assignee_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Item name is required")
	}

	pos, err := h.repo.GetMaxItemPosition(req.ChecklistID)
	if err != nil {
		pos = 0
	}

	item := &models.KanbanChecklistItem{
		ChecklistID: req.ChecklistID,
		Name:        req.Name,
		Position:    pos,
		AssigneeID:  req.AssigneeID,
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(item).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create checklist item: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, item)
}

func (h *KanbanChecklistHandler) UpdateItem(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req struct {
		Name       *string `json:"name"`
		IsChecked  *bool   `json:"is_checked"`
		AssigneeID *uint   `json:"assignee_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	item, err := h.repo.FindItemByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Checklist item not found")
	}

	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.IsChecked != nil {
		item.IsChecked = *req.IsChecked
	}
	if req.AssigneeID != nil {
		item.AssigneeID = req.AssigneeID
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Save(item).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update checklist item: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, item)
}

func (h *KanbanChecklistHandler) DeleteItem(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&models.KanbanChecklistItem{}, id).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete checklist item: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Checklist item deleted"})
}
