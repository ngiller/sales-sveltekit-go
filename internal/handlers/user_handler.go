package handlers

import (
	"backend/config"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) FindAll(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	sortBy := c.Query("sort")
	sortDir := c.Query("order")

	users, total, err := h.repo.FindAll(search, page, limit, sortBy, sortDir)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *UserHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	user, err := h.repo.FindByID(uint(id))
	if err != nil || user == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, user)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Password      string `json:"password"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		PhoneNo       string `json:"phone_no"`
		UserGroupID   *uint  `json:"user_group_id"`
		DepartementID *uint  `json:"departement_id"`
		Enable        bool   `json:"enable"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Password is required")
	}

	hashedPassword, err := config.HashPassword(req.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	// Audit fields
	userID := c.Locals("user_id").(uint)
	phoneNoStr := ""
	if req.PhoneNo != "" {
		phoneNoStr = req.PhoneNo
	}

	user := &models.User{
		Name:          req.Name,
		Email:         req.Email,
		PhoneNo:       &phoneNoStr,
		Password:      hashedPassword,
		UserGroupID:   req.UserGroupID,
		DepartementID: req.DepartementID,
		Enable:        req.Enable,
		UserCreated:   int64(userID),
		UserUpdate:    int64(userID),
	}

	// Use transaction
	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req struct {
		Name          string `json:"name"`
		Email         string `json:"email"`
		PhoneNo       *string `json:"phone_no"`
		UserGroupID   *uint   `json:"user_group_id"`
		DepartementID *uint   `json:"departement_id"`
		PropertyID    *int64  `json:"property_id"`
		Enable        bool    `json:"enable"`
		Inisial       *string `json:"inisial"`
		Password      string  `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := h.repo.FindByID(uint(id))
	if err != nil || user == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	// Only admins (user_group_id == 1) can change role, property, or enable/disable users
	currentUserID := c.Locals("user_id").(uint)
	var currentUser models.User

	hrAdmin := false
	if err := h.repo.GetDB().First(&currentUser, currentUserID).Error; err == nil {
		hrAdmin = currentUser.UserGroupID != nil && *currentUser.UserGroupID == 1
	}

	// Build update map (avoids GORM association issues with Save on preloaded models)
	updates := map[string]interface{}{
		"name":           req.Name,
		"email":          req.Email,
		"phone_no":       req.PhoneNo,
		"departement_id": req.DepartementID,
		"inisial":        req.Inisial,
		"user_update":    int64(currentUserID),
	}

	if hrAdmin {
		updates["user_group_id"] = req.UserGroupID
		updates["property_id"] = req.PropertyID
		updates["enable"] = req.Enable
	}

	if req.Password != "" {
		hashedPassword, err := config.HashPassword(req.Password)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
		}
		updates["password"] = hashedPassword
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, user)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	// Use transaction
	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.User{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete user: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "User successfully deleted"})
}

func (h *UserHandler) UploadAvatar(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded")
	}

	if msg := utils.ValidateFile(file, utils.AllowedImageExts(), 2*1024*1024); msg != "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, msg)
	}

	extension := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), extension)
	savePath := filepath.Join("uploads", "avatars", newFilename)

	if err := c.SaveFile(file, savePath); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save file")
	}

	user, err := h.repo.FindByID(uint(id))
	if err != nil || user == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	fileURL := fmt.Sprintf("%s/%s", c.BaseURL(), filepath.ToSlash(savePath))
	user.Avatar = &fileURL
	if err := h.repo.Update(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user avatar")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Avatar uploaded successfully",
		"url":     fileURL,
	})
}

func (h *UserHandler) UploadSignature(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	file, err := c.FormFile("signature")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded")
	}

	if msg := utils.ValidateFile(file, utils.AllowedImageExts(), 1*1024*1024); msg != "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, msg)
	}

	extension := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), extension)
	savePath := filepath.Join("uploads", "signatures", newFilename)

	if err := c.SaveFile(file, savePath); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save file")
	}

	user, err := h.repo.FindByID(uint(id))
	if err != nil || user == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	fileURL := fmt.Sprintf("%s/%s", c.BaseURL(), filepath.ToSlash(savePath))
	user.Sign = &fileURL
	if err := h.repo.Update(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user signature")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Signature uploaded successfully",
		"url":     fileURL,
	})
}
