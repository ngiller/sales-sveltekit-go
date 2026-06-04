package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CustomerHandler struct {
	repo *repository.CustomerRepository
}

func NewCustomerHandler(repo *repository.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{repo: repo}
}

func (h *CustomerHandler) FindAll(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	sortBy := c.Query("sort")
	sortDir := c.Query("order")

	customers, total, err := h.repo.FindAll(search, page, limit, sortBy, sortDir)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve customers")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"customers": customers,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

func (h *CustomerHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	customer, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Customer not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, customer)
}

func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	var customer models.Customer
	if err := c.BodyParser(&customer); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Audit fields
	userID := c.Locals("user_id").(uint)
	customer.UserCreated = &userID
	customer.UserUpdate = &userID
	customer.UsersID = userID

	// Use a transaction to ensure customer and contacts are saved together
	err := h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&customer).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create customer: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, customer)
}

func (h *CustomerHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var req models.Customer
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	customer, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Customer not found")
	}

	// Update only fields that are likely to be in the form
	customer.CategoryID = req.CategoryID
	customer.Name = req.Name
	customer.Address = req.Address
	customer.Phone = req.Phone
	customer.Email = req.Email
	customer.NPWP = req.NPWP
	customer.Enable = req.Enable
	
	if req.PropertyID != 0 {
		customer.PropertyID = req.PropertyID
	}
	if req.SalesID != nil {
		customer.SalesID = req.SalesID
	}
	if req.AccID != nil {
		customer.AccID = req.AccID
	}

	// Audit fields
	userID := c.Locals("user_id").(uint)
	customer.UserUpdate = &userID

	// Use a transaction to update customer and contacts
	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(customer).Omit("Contacts", "CreatedAt", "UserCreated", "UsersID").Updates(customer).Error; err != nil {
			return err
		}
		
		// Sync contacts association
		for i := range req.Contacts {
			req.Contacts[i].CustomerID = uint(id)
		}
		
		if err := tx.Model(customer).Association("Contacts").Unscoped().Replace(req.Contacts); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		println("ERROR Update Customer:", err.Error())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, customer)
}

func (h *CustomerHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	// Use a transaction to ensure customer and contacts (if not cascaded) are deleted
	err = h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		// Delete contacts first if they are not cascaded (though we added cascade in models)
		if err := tx.Where("customer_id = ?", id).Delete(&models.CustomerContact{}).Error; err != nil {
			return err
		}
		// Delete customer
		if err := tx.Delete(&models.Customer{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete customer: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Customer successfully deleted"})
}
