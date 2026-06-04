package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type KanbanListHandler struct {
	repo      *repository.KanbanListRepository
	boardRepo *repository.KanbanBoardRepository
}

func NewKanbanListHandler(repo *repository.KanbanListRepository, boardRepo *repository.KanbanBoardRepository) *KanbanListHandler {
	return &KanbanListHandler{repo: repo, boardRepo: boardRepo}
}

func (h *KanbanListHandler) FindByBoardID(c *fiber.Ctx) error {
	boardIDStr := c.Query("board_id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "board_id is required")
	}

	lists, err := h.repo.FindByBoardID(uint(boardID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve lists")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, lists)
}

func (h *KanbanListHandler) Create(c *fiber.Ctx) error {
	var req struct {
		BoardID uint   `json:"board_id"`
		Name    string `json:"name"`
		Color   string `json:"color"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "List name is required")
	}

	pos, err := h.repo.GetMaxPosition(req.BoardID)
	if err != nil {
		pos = 0
	}

	list := &models.KanbanList{
		BoardID:  req.BoardID,
		Name:     req.Name,
		Position: pos,
	}
	if req.Color != "" {
		list.Color = &req.Color
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(list).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create list: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, list)
}

func (h *KanbanListHandler) Update(c *fiber.Ctx) error {
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

	list, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "List not found")
	}

	if req.Name != "" {
		list.Name = req.Name
	}
	if req.Color != "" {
		list.Color = &req.Color
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Save(list).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update list: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, list)
}

func (h *KanbanListHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&models.KanbanList{}, id).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete list: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "List successfully deleted"})
}

func (h *KanbanListHandler) Reorder(c *fiber.Ctx) error {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.repo.Reorder(req.IDs); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reorder lists: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Lists reordered"})
}
