package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type KanbanCommentHandler struct {
	repo *repository.KanbanCommentRepository
}

func NewKanbanCommentHandler(repo *repository.KanbanCommentRepository) *KanbanCommentHandler {
	return &KanbanCommentHandler{repo: repo}
}

func (h *KanbanCommentHandler) FindByCardID(c *fiber.Ctx) error {
	cardIDStr := c.Query("card_id")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "card_id is required")
	}

	comments, err := h.repo.FindByCardID(uint(cardID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve comments")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, comments)
}

func (h *KanbanCommentHandler) Create(c *fiber.Ctx) error {
	var req struct {
		CardID  uint   `json:"card_id"`
		Content string `json:"content"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Content == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Comment content is required")
	}

	userID := c.Locals("user_id").(uint)

	comment := &models.KanbanComment{
		CardID:  req.CardID,
		UserID:  &userID,
		Content: req.Content,
	}

	err := h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(comment).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create comment: "+err.Error())
	}

	// Reload with user relation
	comment, _ = h.repo.FindByID(comment.ID)
	return utils.SuccessResponse(c, fiber.StatusCreated, comment)
}

func (h *KanbanCommentHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req struct {
		Content string `json:"content"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	comment, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Comment not found")
	}

	comment.Content = req.Content

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Save(comment).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update comment: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, comment)
}

func (h *KanbanCommentHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&models.KanbanComment{}, id).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete comment: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Comment deleted"})
}
