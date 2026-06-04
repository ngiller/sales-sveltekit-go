package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PolicyHandler struct {
	repo *repository.PolicyRepository
}

func NewPolicyHandler(repo *repository.PolicyRepository) *PolicyHandler {
	return &PolicyHandler{repo: repo}
}

func (h *PolicyHandler) FindAll(c *fiber.Ctx) error {
	policies, err := h.repo.FindAll()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve policies")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, policies)
}

func (h *PolicyHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	policy, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Policy not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, policy)
}

func (h *PolicyHandler) FindByGroupID(c *fiber.Ctx) error {
	groupIDStr := c.Params("groupID")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid group ID format")
	}

	policies, err := h.repo.FindByGroupID(uint(groupID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve policies for group")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, policies)
}

func (h *PolicyHandler) Create(c *fiber.Ctx) error {
	var policy models.GroupPolicy
	if err := c.BodyParser(&policy); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.repo.Create(&policy); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create policy")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, policy)
}

func (h *PolicyHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var reqPolicy models.GroupPolicy
	if err := c.BodyParser(&reqPolicy); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	policy, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Policy not found")
	}

	policy.GroupID = reqPolicy.GroupID
	policy.TargetTableName = reqPolicy.TargetTableName
	policy.TableID = reqPolicy.TableID
	policy.Action = reqPolicy.Action
	policy.PropertyID = reqPolicy.PropertyID

	if err := h.repo.Update(policy); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update policy")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, policy)
}

func (h *PolicyHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete policy")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Policy successfully deleted"})
}
