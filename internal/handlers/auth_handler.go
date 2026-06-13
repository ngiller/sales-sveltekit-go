package handlers

import (
	"backend/config"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
	menuRepo *repository.MenuRepository
}

func NewAuthHandler(userRepo *repository.UserRepository, menuRepo *repository.MenuRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, menuRepo: menuRepo}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, errors)
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		log.Printf("ERROR: FindByEmail: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	if user == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid login credentials provided")
	}

	if !config.CheckPassword(user.Password, req.Password) {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid login credentials provided")
	}

	if !user.Enable {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Account is disabled")
	}

	inisial := ""
	if user.Inisial != nil {
		inisial = *user.Inisial
	}
	token, err := config.GenerateJWT(int64(user.ID), user.Email, inisial, user.UserGroupID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	var menus []models.MenuItem
	if user.UserGroupID != nil {
		menus, _ = h.menuRepo.GetMenusByUserGroupID(*user.UserGroupID)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":             user.ID,
			"name":           user.Name,
			"email":          user.Email,
			"phone_no":       user.PhoneNo,
			"avatar":         user.Avatar,
			"sign":           user.Sign,
			"inisial":        user.Inisial,
			"property_id":    user.PropertyID,
			"user_group_id":  user.UserGroupID,
			"role_name":      user.RoleName,
			"departement_id": user.DepartementID,
			"dept_name":      user.DeptName,
		},
		"menus": menus,
	})
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	var menus []models.MenuItem
	if user.UserGroupID != nil {
		menus, _ = h.menuRepo.GetMenusByUserGroupID(*user.UserGroupID)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"id":             user.ID,
		"name":           user.Name,
		"email":          user.Email,
		"phone_no":       user.PhoneNo,
		"avatar":         user.Avatar,
		"sign":           user.Sign,
		"inisial":        user.Inisial,
		"property_id":    user.PropertyID,
		"user_group_id":  user.UserGroupID,
		"departement_id": user.DepartementID,
		"menus":          menus,
	})
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	loggedInUserID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	idStr := c.Params("id")
	targetUserID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	if uint(targetUserID) != loggedInUserID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only update your own profile")
	}

	user, err := h.userRepo.FindByID(uint(targetUserID))
	if err != nil || user == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	var input struct {
		Name          string  `json:"name"`
		PhoneNo       *string `json:"phone_no"`
		Inisial       *string `json:"inisial"`
		DepartementID *uint  `json:"departement_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.PhoneNo != nil {
		user.PhoneNo = input.PhoneNo
	}
	if input.Inisial != nil {
		user.Inisial = input.Inisial
	}
	if input.DepartementID != nil {
		user.DepartementID = input.DepartementID
	}

	if err := h.userRepo.Update(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"id":             user.ID,
		"name":           user.Name,
		"email":          user.Email,
		"phone_no":       user.PhoneNo,
		"avatar":         user.Avatar,
		"sign":           user.Sign,
		"inisial":        user.Inisial,
		"property_id":    user.PropertyID,
		"user_group_id":  user.UserGroupID,
		"departement_id": user.DepartementID,
	})
}

func (h *AuthHandler) UploadProfileAvatar(c *fiber.Ctx) error {
	loggedInUserID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded: "+err.Error())
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

	user, err := h.userRepo.FindByID(loggedInUserID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	fileURL := fmt.Sprintf("%s/%s", c.BaseURL(), filepath.ToSlash(savePath))
	user.Avatar = &fileURL
	if err := h.userRepo.Update(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update avatar")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Avatar uploaded successfully",
		"url":     fileURL,
	})
}

func (h *AuthHandler) UploadProfileSignature(c *fiber.Ctx) error {
	loggedInUserID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
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

	user, err := h.userRepo.FindByID(loggedInUserID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	fileURL := fmt.Sprintf("%s/%s", c.BaseURL(), filepath.ToSlash(savePath))
	user.Sign = &fileURL
	if err := h.userRepo.Update(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update signature")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Signature uploaded successfully",
		"url":     fileURL,
	})
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	loggedInUserID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	idStr := c.Params("id")
	targetUserID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	if uint(targetUserID) != loggedInUserID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only change your own password")
	}

	user, err := h.userRepo.FindByID(uint(targetUserID))
	if err != nil || user == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	var input struct {
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if input.NewPassword == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "New password is required")
	}

	if len(input.NewPassword) < 8 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "New password must be at least 8 characters")
	}

	hashedPassword, err := config.HashPassword(input.NewPassword)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	user.Password = hashedPassword
	if err := h.userRepo.Update(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to change password")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Password changed successfully",
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
	})
	return utils.SuccessResponse(c, fiber.StatusOK, "Logged out successfully")
}

