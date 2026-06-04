package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type KanbanCardHandler struct {
	repo    *repository.KanbanCardRepository
	listRepo *repository.KanbanListRepository
}

func NewKanbanCardHandler(repo *repository.KanbanCardRepository, listRepo *repository.KanbanListRepository) *KanbanCardHandler {
	return &KanbanCardHandler{repo: repo, listRepo: listRepo}
}

func (h *KanbanCardHandler) FindByListID(c *fiber.Ctx) error {
	listIDStr := c.Query("list_id")
	listID, err := strconv.ParseUint(listIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "list_id is required")
	}

	cards, err := h.repo.FindByListID(uint(listID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve cards")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, cards)
}

func (h *KanbanCardHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	card, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Card not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, card)
}

func (h *KanbanCardHandler) Create(c *fiber.Ctx) error {
	var req struct {
		ListID      uint       `json:"list_id"`
		BoardID     uint       `json:"board_id"`
		Title       string     `json:"title"`
		Description *string    `json:"description"`
		DueDate     *time.Time `json:"due_date"`
		StartDate   *time.Time `json:"start_date"`
		MemberIDs   []uint     `json:"member_ids"`
		LabelIDs    []uint     `json:"label_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Title == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Card title is required")
	}

	userID := c.Locals("user_id").(uint)

	pos, err := h.repo.GetMaxPosition(req.ListID)
	if err != nil {
		pos = 0
	}

	card := &models.KanbanCard{
		ListID:      req.ListID,
		BoardID:     req.BoardID,
		Title:       req.Title,
		Description: req.Description,
		Position:    pos,
		DueDate:     req.DueDate,
		StartDate:   req.StartDate,
		UserCreated: &userID,
		UserUpdated: &userID,
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(card).Error; err != nil {
			return err
		}
		if len(req.MemberIDs) > 0 {
			if err := tx.Model(card).Association("Members").Replace(req.MemberIDs); err != nil {
				return err
			}
		}
		if len(req.LabelIDs) > 0 {
			if err := tx.Model(card).Association("Labels").Replace(req.LabelIDs); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create card: "+err.Error())
	}

	// Reload with relations
	card, _ = h.repo.FindByID(card.ID)
	return utils.SuccessResponse(c, fiber.StatusCreated, card)
}

func (h *KanbanCardHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req struct {
		Title       string     `json:"title"`
		Description *string    `json:"description"`
		DueDate     *time.Time `json:"due_date"`
		StartDate   *time.Time `json:"start_date"`
		CoverImage  *string    `json:"cover_image"`
		IsArchived  *bool      `json:"is_archived"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	card, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Card not found")
	}

	if req.Title != "" {
		card.Title = req.Title
	}
	if req.Description != nil {
		card.Description = req.Description
	}
	if req.DueDate != nil {
		card.DueDate = req.DueDate
	}
	if req.StartDate != nil {
		card.StartDate = req.StartDate
	}
	if req.CoverImage != nil {
		card.CoverImage = req.CoverImage
	}
	if req.IsArchived != nil {
		card.IsArchived = *req.IsArchived
	}

	userID := c.Locals("user_id").(uint)
	card.UserUpdated = &userID

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Save(card).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update card: "+err.Error())
	}

	card, _ = h.repo.FindByID(card.ID)
	return utils.SuccessResponse(c, fiber.StatusOK, card)
}

func (h *KanbanCardHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&models.KanbanCard{}, id).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete card: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Card successfully deleted"})
}

func (h *KanbanCardHandler) Move(c *fiber.Ctx) error {
	var req struct {
		CardID    uint `json:"card_id"`
		ToListID  uint `json:"to_list_id"`
		Position  int  `json:"position"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	err := h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return h.repo.MoveCard(req.CardID, req.ToListID, req.Position)
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to move card: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Card moved successfully"})
}

func (h *KanbanCardHandler) Reorder(c *fiber.Ctx) error {
	var req struct {
		IDs    []uint `json:"ids"`
		ListID uint   `json:"list_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.repo.Reorder(req.IDs, req.ListID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reorder cards: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Cards reordered"})
}

func (h *KanbanCardHandler) SyncMembers(c *fiber.Ctx) error {
	cardIDStr := c.Params("id")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req struct {
		MemberIDs []uint `json:"member_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.repo.SyncMembers(uint(cardID), req.MemberIDs); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to sync members: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Members updated"})
}

func (h *KanbanCardHandler) SyncLabels(c *fiber.Ctx) error {
	cardIDStr := c.Params("id")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req struct {
		LabelIDs []uint `json:"label_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.repo.SyncLabels(uint(cardID), req.LabelIDs); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to sync labels: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Labels updated"})
}
