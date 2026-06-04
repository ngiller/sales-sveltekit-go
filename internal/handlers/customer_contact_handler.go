package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CustomerContactHandler struct {
	repo *repository.CustomerContactRepository
}

func NewCustomerContactHandler(repo *repository.CustomerContactRepository) *CustomerContactHandler {
	return &CustomerContactHandler{repo: repo}
}

func (h *CustomerContactHandler) FindAllByCustomer(c *fiber.Ctx) error {
	customerIDStr := c.Query("customer_id")
	if customerIDStr == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "customer_id is required")
	}
	
	customerID, err := strconv.ParseUint(customerIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid customer_id format")
	}

	contacts, err := h.repo.FindAllByCustomerID(uint(customerID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve contacts")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, contacts)
}

func (h *CustomerContactHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	contact, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Contact not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, contact)
}

func (h *CustomerContactHandler) Create(c *fiber.Ctx) error {
	var contact models.CustomerContact
	if err := c.BodyParser(&contact); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.repo.Create(&contact); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create contact")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, contact)
}

func (h *CustomerContactHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req models.CustomerContact
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	contact, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Contact not found")
	}

	contact.Name = req.Name
	contact.Phone = req.Phone
	contact.Email = req.Email
	contact.Position = req.Position

	if err := h.repo.Update(contact); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update contact")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, contact)
}

func (h *CustomerContactHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete contact")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Contact successfully deleted"})
}
