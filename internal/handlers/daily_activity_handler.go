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

type DailyActivityHandler struct {
	repo *repository.DailyActivityRepository
}

func NewDailyActivityHandler(repo *repository.DailyActivityRepository) *DailyActivityHandler {
	return &DailyActivityHandler{repo: repo}
}

func getIndonesianDay(t time.Time) string {
	days := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	return days[t.Weekday()]
}

func (h *DailyActivityHandler) FindAll(c *fiber.Ctx) error {
	currentUserID := c.Locals("user_id").(uint)

	// Check if user is Admin (group_id = 1)
	var user models.User
	if err := h.repo.GetDB().Select("user_group_id").First(&user, currentUserID).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}
	isAdmin := user.UserGroupID != nil && *user.UserGroupID == 1

	var userIDFilter uint
	if isAdmin {
		if uFilter := c.Query("user_id"); uFilter != "" {
			val, _ := strconv.ParseUint(uFilter, 10, 32)
			userIDFilter = uint(val)
		}
	} else {
		// Non-admin can only see their own activities
		userIDFilter = currentUserID
	}

	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	sortBy := c.Query("sort")
	sortDir := c.Query("order")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	items, total, err := h.repo.FindAll(userIDFilter, search, page, limit, sortBy, sortDir, startDate, endDate)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve daily activities: "+err.Error())
	}

	// Map day name to Hari
	for i := range items {
		items[i].Hari = getIndonesianDay(items[i].ActivityDate)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *DailyActivityHandler) FindByID(c *fiber.Ctx) error {
	currentUserID := c.Locals("user_id").(uint)

	// Check admin
	var user models.User
	if err := h.repo.GetDB().Select("user_group_id").First(&user, currentUserID).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}
	isAdmin := user.UserGroupID != nil && *user.UserGroupID == 1

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	item, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Daily activity not found")
	}

	// Permission check
	if item.UserID != currentUserID && !isAdmin {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Access Denied: You cannot view other user's activity")
	}

	item.Hari = getIndonesianDay(item.ActivityDate)

	return utils.SuccessResponse(c, fiber.StatusOK, item)
}

type DailyActivityInput struct {
	ID           uint    `json:"id"`
	UserID       uint    `json:"user_id"`
	ActivityDate *string `json:"activity_date"`
	Activity     string  `json:"activity"`
	Process      *string `json:"process"`
	Issues       *string `json:"issues"`
	Result       *string `json:"result"`
	Notes        *string `json:"notes"`
	PropertyID   *uint   `json:"property_id"`
}

func (h *DailyActivityHandler) Create(c *fiber.Ctx) error {
	currentUserID := c.Locals("user_id").(uint)

	var input DailyActivityInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	var item models.DailyActivity
	item.ID = input.ID
	if input.ActivityDate == nil || *input.ActivityDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Activity date is required")
	}
	t := parseDate(input.ActivityDate)
	if t == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid activity date format")
	}
	item.ActivityDate = *t
	item.Activity = input.Activity
	item.Process = input.Process
	item.Issues = input.Issues
	item.Result = input.Result
	item.Notes = input.Notes
	item.PropertyID = input.PropertyID

	// Force current user ID (Standard users can only record their own activities)
	// Even admins, when logging, will log under their own user ID.
	item.UserID = currentUserID

	err := h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create daily activity: "+err.Error())
	}

	item.Hari = getIndonesianDay(item.ActivityDate)

	return utils.SuccessResponse(c, fiber.StatusCreated, item)
}

func (h *DailyActivityHandler) Update(c *fiber.Ctx) error {
	currentUserID := c.Locals("user_id").(uint)

	// Check admin
	var user models.User
	if err := h.repo.GetDB().Select("user_group_id").First(&user, currentUserID).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}
	isAdmin := user.UserGroupID != nil && *user.UserGroupID == 1

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var input DailyActivityInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	item, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Daily activity not found")
	}

	// Permission check
	if item.UserID != currentUserID && !isAdmin {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Access Denied: You cannot update other user's activity")
	}

	// Update fields
	if input.ActivityDate != nil {
		t := parseDate(input.ActivityDate)
		if t != nil {
			item.ActivityDate = *t
		} else {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid activity date format")
		}
	}
	item.Activity = input.Activity
	item.Process = input.Process
	item.Issues = input.Issues
	item.Result = input.Result
	item.Notes = input.Notes
	if input.PropertyID != nil {
		item.PropertyID = input.PropertyID
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(item).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update daily activity: "+err.Error())
	}

	item.Hari = getIndonesianDay(item.ActivityDate)

	return utils.SuccessResponse(c, fiber.StatusOK, item)
}

func (h *DailyActivityHandler) Delete(c *fiber.Ctx) error {
	currentUserID := c.Locals("user_id").(uint)

	// Check admin
	var user models.User
	if err := h.repo.GetDB().Select("user_group_id").First(&user, currentUserID).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}
	isAdmin := user.UserGroupID != nil && *user.UserGroupID == 1

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	item, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Daily activity not found")
	}

	// Permission check
	if item.UserID != currentUserID && !isAdmin {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Access Denied: You cannot delete other user's activity")
	}

	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.DailyActivity{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete daily activity: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Daily activity successfully deleted"})
}
