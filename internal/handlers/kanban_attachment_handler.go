package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KanbanAttachmentHandler struct {
	repo *repository.KanbanAttachmentRepository
}

func NewKanbanAttachmentHandler(repo *repository.KanbanAttachmentRepository) *KanbanAttachmentHandler {
	return &KanbanAttachmentHandler{repo: repo}
}

func (h *KanbanAttachmentHandler) FindByCardID(c *fiber.Ctx) error {
	cardIDStr := c.Query("card_id")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "card_id is required")
	}

	attachments, err := h.repo.FindByCardID(uint(cardID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve attachments")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, attachments)
}

func (h *KanbanAttachmentHandler) Create(c *fiber.Ctx) error {
	cardIDStr := c.FormValue("card_id")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "card_id is required")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded")
	}

	extension := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), extension)
	saveDir := filepath.Join("uploads", "kanban")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create upload directory")
	}
	savePath := filepath.Join(saveDir, newFilename)

	if err := c.SaveFile(file, savePath); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save file")
	}

	userID := c.Locals("user_id").(uint)

	attachment := &models.KanbanAttachment{
		CardID:     uint(cardID),
		FileName:   file.Filename,
		FilePath:   filepath.ToSlash(savePath),
		FileSize:   file.Size,
		MimeType:   file.Header.Get("Content-Type"),
		UploadedBy: &userID,
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(attachment).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save attachment: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, attachment)
}

func (h *KanbanAttachmentHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	attachment, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Attachment not found")
	}

	// Delete file from disk
	fullPath := filepath.Join(".", attachment.FilePath)
	if err := removeFile(fullPath); err != nil {
		fmt.Printf("Warning: could not delete file %s: %v\n", fullPath, err)
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&models.KanbanAttachment{}, id).Error
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete attachment: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Attachment deleted"})
}

func removeFile(path string) error {
	return os.Remove(path)
}
