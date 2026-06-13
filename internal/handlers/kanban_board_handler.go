package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type KanbanBoardHandler struct {
	boardRepo *repository.KanbanBoardRepository
	listRepo  *repository.KanbanListRepository
	cardRepo  *repository.KanbanCardRepository
}

func NewKanbanBoardHandler(boardRepo *repository.KanbanBoardRepository, listRepo *repository.KanbanListRepository, cardRepo *repository.KanbanCardRepository) *KanbanBoardHandler {
	return &KanbanBoardHandler{boardRepo: boardRepo, listRepo: listRepo, cardRepo: cardRepo}
}

func (h *KanbanBoardHandler) FindAll(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	boards, total, err := h.boardRepo.FindAll(search, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve boards")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"boards": boards,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

func (h *KanbanBoardHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	board, err := h.boardRepo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, 		"Board not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, board)
}

func (h *KanbanBoardHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Background  *string `json:"background"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Board name is required")
	}

	userID := c.Locals("user_id").(uint)

	board := &models.KanbanBoard{
		Name:        req.Name,
		Description: req.Description,
		Background:  req.Background,
		UserCreated: &userID,
	}

	err := h.boardRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(board).Error; err != nil {
			return err
		}
		defaultLists := []string{"To Do", "In Progress", "Done"}
		for i, name := range defaultLists {
			list := &models.KanbanList{
				BoardID:  board.ID,
				Name:     name,
				Position: i,
			}
			if err := tx.Create(list).Error; err != nil {
				return err
			}
		}
		defaultLabels := []struct{ Name, Color string }{
			{"Bug", "#ef4444"},
			{"Feature", "#3b82f6"},
			{"Improvement", "#10b981"},
			{"Question", "#f59e0b"},
			{"Urgent", "#ec4899"},
		}
		for _, l := range defaultLabels {
			label := &models.KanbanLabel{
				BoardID: board.ID,
				Name:    l.Name,
				Color:   l.Color,
			}
			if err := tx.Create(label).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create board: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, board)
}

func (h *KanbanBoardHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Background  *string `json:"background"`
		IsArchived  *bool   `json:"is_archived"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	board, err := h.boardRepo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Board not found")
	}

	if req.Name != "" {
		board.Name = req.Name
	}
	if req.Description != nil {
		board.Description = req.Description
	}
	if req.Background != nil {
		board.Background = req.Background
	}
	if req.IsArchived != nil {
		board.IsArchived = *req.IsArchived
	}

	err = h.boardRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Save(board).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update board: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, board)
}

func (h *KanbanBoardHandler) ListUsers(c *fiber.Ctx) error {
	var users []struct {
		ID      uint    `json:"id"`
		Name    string  `json:"name"`
		Inisial *string `json:"inisial"`
	}
	if err := h.boardRepo.GetDB().Model(&models.User{}).Select("id, name, inisial").Where("enable = ?", 1).Find(&users).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"users": users,
	})
}

func (h *KanbanBoardHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	err = h.boardRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&models.KanbanBoard{}, id).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete board: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Board successfully deleted"})
}
