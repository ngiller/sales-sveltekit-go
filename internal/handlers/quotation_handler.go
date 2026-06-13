package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type QuotationHandler struct {
	repo *repository.QuotationRepository
}

func NewQuotationHandler(repo *repository.QuotationRepository) *QuotationHandler {
	return &QuotationHandler{repo: repo}
}

func (h *QuotationHandler) fetchDropdownData() (*fiber.Map, error) {
	db := h.repo.GetDB()

	var users []models.User
	var customers []models.Customer
	var paymentTerms []models.PaymentTerm
	var projectLevels []models.ProjectLevel
	var progress []models.QuotationProgress
	var statuses []models.QuotationStatus
	var units []models.Unit
	var categories []models.CustomerCategory

	if err := db.Where("enable = ?", "1").Find(&users).Error; err != nil {
		return nil, err
	}
	if err := db.Where("enable = ?", "1").Find(&customers).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&paymentTerms).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&projectLevels).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&progress).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&statuses).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&units).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return &fiber.Map{
		"users":          users,
		"customers":      customers,
		"payment_terms":  paymentTerms,
		"project_levels": projectLevels,
		"progress":       progress,
		"statuses":       statuses,
		"units":          units,
		"categories":     categories,
	}, nil
}

func (h *QuotationHandler) New(c *fiber.Ctx) error {
	dropdownData, err := h.fetchDropdownData()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve dropdown data: "+err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusOK, *dropdownData)
}

func (h *QuotationHandler) fetchReportSalesPersons(userID uint) []fiber.Map {
	canViewAll := h.userCanViewAllQuotations(userID)

	var users []models.User
	query := h.repo.GetDB().Model(&models.User{}).Select("id, name, inisial")
	if !canViewAll {
		query = query.Where("id = ?", userID)
	}
	query.Find(&users)

	result := make([]fiber.Map, len(users))
	for i, u := range users {
		result[i] = fiber.Map{
			"id":      u.ID,
			"name":    u.Name,
			"inisial": u.Inisial,
		}
	}
	return result
}

func (h *QuotationHandler) userCanViewAllQuotations(userID uint) bool {
	var user models.User
	if err := h.repo.GetDB().Select("user_group_id").First(&user, userID).Error; err != nil {
		return false
	}

	// Admin bypass
	if user.UserGroupID != nil && *user.UserGroupID == 1 {
		return true
	}

	// Check for view_all policy on quotations endpoint
	var table models.MasterTableAccess
	if err := h.repo.GetDB().Where("endpoint = ?", "quotations").First(&table).Error; err != nil {
		return false
	}

	var count int64
	h.repo.GetDB().Model(&models.GroupPolicy{}).
		Where("group_id = ? AND action = ? AND (table_name = ? OR table_id = ?)",
			*user.UserGroupID, "view_all", table.Name, table.ID).
		Count(&count)
	return count > 0
}

func (h *QuotationHandler) FindAll(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	search := c.Query("search")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	qType := c.Query("quotation_type")
	status := c.Query("status")
	userCreated := c.Query("sales_id")
	if userCreated == "" {
		userCreated = c.Query("user_created") // fallback
	}
	progress := c.Query("progress")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	sortBy := c.Query("sort")
	sortDir := c.Query("order")

	canViewAll := h.userCanViewAllQuotations(userID)
	if !canViewAll {
		userCreated = strconv.FormatUint(uint64(userID), 10)
	}

	items, total, err := h.repo.FindAll(search, fromDate, toDate, qType, status, userCreated, progress, page, limit, sortBy, sortDir)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve quotations: "+err.Error())
	}


	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items":        items,
		"total":        total,
		"page":         page,
		"limit":        limit,
		"can_view_all": canViewAll,
	})
}

func (h *QuotationHandler) FindByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Quotation ID is required")
	}

	quotation, err := h.repo.FindByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Quotation not found: "+err.Error())
	}

	// The repository FindByID already preloads the default revision's details
	// but we still fetch them separately for the frontend's legacy structure if needed
	// Find the default revID from revisions

	dropdownData, err := h.fetchDropdownData()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve dropdown data: "+err.Error())
	}

	details := quotation.Details
	subdetails := quotation.Subdetails
	revisions := quotation.Revisions

	result := fiber.Map{
		"quotation":  quotation,
		"details":    details,
		"subdetails": subdetails,
		"revisions":  revisions,
	}
	for k, v := range *dropdownData {
		result[k] = v
	}

	return utils.SuccessResponse(c, fiber.StatusOK, result)
}

type QuotationCreateUpdateInput struct {
	PropertyID    int           `json:"property_id"`
	QuotationType int           `json:"quotation_type"`
	QuotationDate *string       `json:"quotation_date"`
	CustomerID    uint          `json:"customer_id"`
	ContactID     *uint         `json:"contact_id"`
	Subject       *string       `json:"subject"`
	Total         *float64      `json:"total"`
	Tax           *float64      `json:"tax"`
	TaxValue      float64       `json:"tax_value"`
	PPh           float64       `json:"pph"`
	PPhValue      float64       `json:"pph_value"`
	Disc          *float64      `json:"disc"`
	GrandTotal    *float64      `json:"grand_total"`
	HppTotal      *float64      `json:"hpp_total"`
	Profit        *float64      `json:"profit"`
	ProfitValue   *float64      `json:"profit_value"`
	PaymentTermID *uint         `json:"payment_term_id"`
	ValidUntil    *string       `json:"valid_until"`
	Commision     float64       `json:"commision"`
	Notes         *string       `json:"notes"`
	Status        int           `json:"status"`
	ProgressID    uint          `json:"progress_id"`
	FollowupBy    *uint         `json:"followup_by"`
	FollowupDate  *string       `json:"followup_date"`
	NextFollowup  *string       `json:"next_followup"`
	Folder        *string       `json:"folder"`
	LevelID       *uint         `json:"level_id"`
	PriorityID    *uint         `json:"priority_id"`
	ProjectStart  *string       `json:"project_start"`
	ProjectEnd    *string       `json:"project_end"`
	PoNo          *string       `json:"po_no"`
	PoDate        *string       `json:"po_date"`
	PoFile        *string       `json:"po_file"`
	PoAssignTo    *string       `json:"po_assign_to"`
	SalesID       *uint         `json:"sales_id"`
	UserCreated   *uint         `json:"user_created"`
	Details       []DetailInput `json:"details"`
}

type DetailInput struct {
	Line        int              `json:"line"`
	No          *int             `json:"no"`
	ProductID   *uint            `json:"product_id"`
	PartNo      *string          `json:"part_no"`
	ProductType int              `json:"product_type"`
	Description *string          `json:"description"`
	Qty         *float64         `json:"qty"`
	UnitID      *uint            `json:"unit_id"`
	Price       float64          `json:"price"`
	Total       float64          `json:"total"`
	OtherCost   float64          `json:"other_cost"`
	Hpp         *float64         `json:"hpp"`
	HppTotal    *float64         `json:"hpp_total"`
	Children    []SubdetailInput `json:"children"`
}

type SubdetailInput struct {
	Line        int      `json:"line"`
	Subline     int      `json:"subline"`
	No          *int     `json:"no"`
	ProductID   *uint    `json:"product_id"`
	PartNo      *string  `json:"part_no"`
	ProductType int      `json:"product_type"`
	Description *string  `json:"description"`
	Qty         *float64 `json:"qty"`
	UnitID      *uint    `json:"unit_id"`
	Price       float64  `json:"price"`
	Total       float64  `json:"total"`
	OtherCost   float64  `json:"other_cost"`
	Hpp         *float64 `json:"hpp"`
	HppTotal    *float64 `json:"hpp_total"`
}

func getUID(c *fiber.Ctx) uint {
	userID := c.Locals("user_id")
	if userID == nil {
		return 0
	}
	switch v := userID.(type) {
	case uint:
		return v
	case int:
		return uint(v)
	case float64:
		return uint(v)
	default:
		return 0
	}
}

func parseDate(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil
	}
	// MySQL DATE range is 1000-01-01 to 9999-12-31
	if t.Year() < 1000 || t.Year() > 9999 {
		return nil
	}
	return &t
}

func getFloatFromPtr(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func (h *QuotationHandler) Create(c *fiber.Ctx) error {
	defer func() {
		recover()
	}()
	uid := getUID(c)
	userInisial := c.Locals("inisial")
	if uid == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var input QuotationCreateUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if input.CustomerID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Customer is required")
	}
	if input.QuotationType == 0 {
		input.QuotationType = 1
	}

	newID := repository.GenerateID()
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	ym := now.Format("200601")

	counter, err := h.repo.GetCounter(input.PropertyID, "quotation", input.QuotationType, ym)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate counter")
	}

	userInisialStr := "USR"
	if inisial, ok := userInisial.(string); ok && inisial != "" {
		userInisialStr = inisial
	}
	quotationID := repository.GenerateQuotationID(input.QuotationType, counter, month, year, userInisialStr)

	quotation := models.Quotation{
		ID:            newID,
		PropertyID:    input.PropertyID,
		QuotationType: input.QuotationType,
		QuotationID:   quotationID,
		QuotationDate: parseDate(input.QuotationDate),
		CustomerID:    input.CustomerID,
		ContactID:     input.ContactID,
		Subject:       input.Subject,
		Total:         input.Total,
		Tax:           input.Tax,
		TaxValue:      input.TaxValue,
		PPh:           input.PPh,
		PPhValue:      input.PPhValue,
		Disc:          input.Disc,
		GrandTotal:    input.GrandTotal,
		HppTotal:      input.HppTotal,
		Profit:        input.Profit,
		ProfitValue:   input.ProfitValue,
		PaymentTermID: input.PaymentTermID,
		ValidUntil:    parseDate(input.ValidUntil),
		Commision:     input.Commision,
		Notes:         input.Notes,
		Status:        input.Status,
		ProgressID:    input.ProgressID,
		FollowupBy:    input.FollowupBy,
		FollowupDate:  parseDate(input.FollowupDate),
		NextFollowup:  parseDate(input.NextFollowup),
		Folder:        input.Folder,
		LevelID:       input.LevelID,
		PriorityID:    input.PriorityID,
		ProjectStart:  parseDate(input.ProjectStart),
		ProjectEnd:    parseDate(input.ProjectEnd),
		PoNo:          input.PoNo,
		PoDate:        parseDate(input.PoDate),
		PoFile:        input.PoFile,
		PoAssignTo:    input.PoAssignTo,
		SalesID:       input.SalesID,
	}

	quotation.SalesID = input.SalesID
	quotation.UserCreated = &uid
	quotation.UserUpdate = &uid
	nowTime := time.Now()
	quotation.CreatedAt = nowTime
	quotation.UpdatedAt = nowTime

	validUntilDate := now.AddDate(0, 0, 7)
	quotation.ValidUntil = &validUntilDate
	nextFollowupDate := now.AddDate(0, 0, 3)
	quotation.NextFollowup = &nextFollowupDate

	if input.Status == 0 {
		input.Status = 1
		quotation.Status = 1
	}
	if input.ProgressID == 0 {
		input.ProgressID = 2
		quotation.ProgressID = 2
	}

	subdetails := []models.QuotationSubdetail{}
	for _, d := range input.Details {
		for _, s := range d.Children {
			subdetail := models.QuotationSubdetail{
				ID:          newID,
				RevID:       0,
				Line:        d.Line,
				Subline:     s.Subline,
				No:          s.No,
				ProductID:   s.ProductID,
				PartNo:      s.PartNo,
				ProductType: s.ProductType,
				Description: s.Description,
				Qty:         s.Qty,
				UnitID:      s.UnitID,
				Price:       s.Price,
				Total:       s.Total,
				OtherCost:   s.OtherCost,
				Hpp:         s.Hpp,
				HppTotal:    s.HppTotal,
			}
			subdetails = append(subdetails, subdetail)
		}
	}

	var details []models.QuotationDetail
	for _, d := range input.Details {
		detail := models.QuotationDetail{
			ID:          newID,
			RevID:       0,
			Line:        d.Line,
			No:          d.No,
			ProductID:   d.ProductID,
			PartNo:      d.PartNo,
			ProductType: d.ProductType,
			Description: d.Description,
			Qty:         d.Qty,
			UnitID:      d.UnitID,
			Price:       d.Price,
			Total:       d.Total,
			OtherCost:   d.OtherCost,
			Hpp:         d.Hpp,
			HppTotal:    d.HppTotal,
		}
		details = append(details, detail)
	}

	master := &models.QuotationMaster{
		ID:            newID,
		RevID:         0,
		QuotationDate: parseDate(input.QuotationDate),
		Subject:       input.Subject,
		Total:         input.Total,
		Disc:          input.Disc,
		Tax:           input.Tax,
		TaxValue:      input.TaxValue,
		PPh:           input.PPh,
		PPhValue:      input.PPhValue,
		GrandTotal:    input.GrandTotal,
		HppTotal:      input.HppTotal,
		Profit:        getFloatFromPtr(input.Profit),
		ProfitValue:   input.ProfitValue,
		PaymentTermID: input.PaymentTermID,
		ValidUntil:    parseDate(input.ValidUntil),
		Commision:     input.Commision,
		Notes:         input.Notes,
		SalesID:       input.SalesID,
		UserCreated:   &uid,
		UserUpdate:    &uid,
		CreatedAt:     nowTime,
		UpdatedAt:     nowTime,
		DefaultQuot:   true,
		LevelID:       input.LevelID,
		PriorityID:    input.PriorityID,
		ProjectStart:  parseDate(input.ProjectStart),
		ProjectEnd:    parseDate(input.ProjectEnd),
		PoNo:          input.PoNo,
		PoDate:        parseDate(input.PoDate),
		PoFile:        input.PoFile,
		PoAssignTo:    input.PoAssignTo,
	}


	err = h.repo.Create(&quotation, details, subdetails, master)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create quotation: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, fiber.Map{
		"quotation":    quotation,
		"quotation_id": quotationID,
	})
}

func (h *QuotationHandler) Update(c *fiber.Ctx) error {
	uid := getUID(c)
	if uid == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	id := c.Params("id")
	if id == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Quotation ID is required")
	}

	existing, err := h.repo.FindByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Quotation not found")
	}

	var user models.User
	if err := h.repo.GetDB().First(&user, uid).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	isAdmin := user.UserGroupID != nil && *user.UserGroupID == 1
	if !isAdmin {
		if existing.UserCreated != nil && *existing.UserCreated != uid {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not authorized to modify this quotation")
		}
	}

	var input QuotationCreateUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	existing.PropertyID = input.PropertyID
	existing.QuotationType = input.QuotationType
	existing.QuotationDate = parseDate(input.QuotationDate)
	existing.CustomerID = input.CustomerID
	existing.ContactID = input.ContactID
	existing.Subject = input.Subject
	existing.Total = input.Total
	existing.Tax = input.Tax
	existing.TaxValue = input.TaxValue
	existing.PPh = input.PPh
	existing.PPhValue = input.PPhValue
	existing.Disc = input.Disc
	existing.GrandTotal = input.GrandTotal
	existing.HppTotal = input.HppTotal
	existing.Profit = input.Profit
	existing.ProfitValue = input.ProfitValue
	existing.PaymentTermID = input.PaymentTermID
	existing.ValidUntil = parseDate(input.ValidUntil)
	existing.Commision = input.Commision
	existing.Notes = input.Notes
	existing.Status = input.Status
	existing.ProgressID = input.ProgressID
	existing.FollowupBy = input.FollowupBy
	existing.FollowupDate = parseDate(input.FollowupDate)
	existing.NextFollowup = parseDate(input.NextFollowup)
	existing.Folder = input.Folder
	existing.LevelID = input.LevelID
	existing.PriorityID = input.PriorityID
	existing.ProjectStart = parseDate(input.ProjectStart)
	existing.ProjectEnd = parseDate(input.ProjectEnd)
	existing.PoNo = input.PoNo
	existing.PoDate = parseDate(input.PoDate)
	existing.PoFile = input.PoFile
	existing.PoAssignTo = input.PoAssignTo
	existing.SalesID = input.SalesID
	if input.UserCreated != nil {
		existing.UserCreated = input.UserCreated
	}
	existing.UserUpdate = &uid
	existing.UpdatedAt = time.Now()

	var details []models.QuotationDetail
	for _, d := range input.Details {
		detail := models.QuotationDetail{
			ID:          id,
			RevID:       0,
			Line:        d.Line,
			No:          d.No,
			ProductID:   d.ProductID,
			PartNo:      d.PartNo,
			ProductType: d.ProductType,
			Description: d.Description,
			Qty:         d.Qty,
			UnitID:      d.UnitID,
			Price:       d.Price,
			Total:       d.Total,
			OtherCost:   d.OtherCost,
			Hpp:         d.Hpp,
			HppTotal:    d.HppTotal,
		}
		details = append(details, detail)
	}

	subdetails := []models.QuotationSubdetail{}
	for _, d := range input.Details {
		for _, s := range d.Children {
			subdetail := models.QuotationSubdetail{
				ID:          id,
				RevID:       0,
				Line:        d.Line,
				Subline:     s.Subline,
				No:          s.No,
				ProductID:   s.ProductID,
				PartNo:      s.PartNo,
				ProductType: s.ProductType,
				Description: s.Description,
				Qty:         s.Qty,
				UnitID:      s.UnitID,
				Price:       s.Price,
				Total:       s.Total,
				OtherCost:   s.OtherCost,
				Hpp:         s.Hpp,
				HppTotal:    s.HppTotal,
			}
			subdetails = append(subdetails, subdetail)
		}
	}

	master := &models.QuotationMaster{
		ID:            id,
		RevID:         -1, // Signal to repository to auto-increment
		QuotationDate: parseDate(input.QuotationDate),
		Subject:       input.Subject,
		Total:         input.Total,
		Disc:          input.Disc,
		Tax:           input.Tax,
		TaxValue:      input.TaxValue,
		PPh:           input.PPh,
		PPhValue:      input.PPhValue,
		GrandTotal:    input.GrandTotal,
		HppTotal:      input.HppTotal,
		Profit:        getFloatFromPtr(input.Profit),
		ProfitValue:   input.ProfitValue,
		PaymentTermID: input.PaymentTermID,
		ValidUntil:    parseDate(input.ValidUntil),
		Commision:     input.Commision,
		Notes:         input.Notes,
		SalesID:       input.SalesID,
		UserCreated:   existing.UserCreated,
		UserUpdate:    &uid,
		UpdatedAt:     time.Now(),
		DefaultQuot:   true,
		LevelID:       input.LevelID,
		PriorityID:    input.PriorityID,
		ProjectStart:  parseDate(input.ProjectStart),
		ProjectEnd:    parseDate(input.ProjectEnd),
		PoNo:          input.PoNo,
		PoDate:        parseDate(input.PoDate),
		PoFile:        input.PoFile,
		PoAssignTo:    input.PoAssignTo,
	}

	err = h.repo.Update(existing, details, subdetails, master)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update quotation: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, existing)
}


func (h *QuotationHandler) Delete(c *fiber.Ctx) error {
	uid := getUID(c)
	if uid == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	id := c.Params("id")
	if id == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Quotation ID is required")
	}

	existing, err := h.repo.FindByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Quotation not found")
	}

	var user models.User
	if err := h.repo.GetDB().First(&user, uid).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	isAdmin := user.UserGroupID != nil && *user.UserGroupID == 1
	if !isAdmin {
		if existing.UserCreated != nil && *existing.UserCreated != uid {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not authorized to delete this quotation")
		}
	}

	err = h.repo.Delete(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete quotation: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Quotation deleted successfully"})
}

func (h *QuotationHandler) CreateRevision(c *fiber.Ctx) error {
	uid := getUID(c)
	userInisial := c.Locals("inisial")
	if uid == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	oldID := c.Params("id")
	if oldID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Quotation ID is required")
	}

	oldQ, err := h.repo.FindByID(oldID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Quotation not found")
	}

	var user models.User
	if err := h.repo.GetDB().First(&user, uid).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	isAdmin := user.UserGroupID != nil && *user.UserGroupID == 1
	if !isAdmin {
		if oldQ.UserCreated != nil && *oldQ.UserCreated != uid {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not authorized to create a revision of this quotation")
		}
	}

	newID := repository.GenerateID()
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	ym := now.Format("200601")

	counter, err := h.repo.GetCounter(oldQ.PropertyID, "quotation", oldQ.QuotationType, ym)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate counter")
	}

	userInisialStr := "USR"
	if inisial, ok := userInisial.(string); ok && inisial != "" {
		userInisialStr = inisial
	}
	newQuotationID := repository.GenerateQuotationID(oldQ.QuotationType, counter, month, year, userInisialStr)

	err = h.repo.CreateRevision(oldID, newID, newQuotationID, now, uid, userInisialStr, "[REVISION] ")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create revision: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, fiber.Map{
		"message":      "Revision created successfully",
		"new_id":       newID,
		"quotation_id": newQuotationID,
	})
}

func (h *QuotationHandler) Duplicate(c *fiber.Ctx) error {
	uid := getUID(c)
	userInisial := c.Locals("inisial")
	if uid == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	oldID := c.Params("id")
	if oldID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Quotation ID is required")
	}

	oldQ, err := h.repo.FindByID(oldID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Quotation not found")
	}

	var user models.User
	if err := h.repo.GetDB().First(&user, uid).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	isAdmin := user.UserGroupID != nil && *user.UserGroupID == 1
	if !isAdmin {
		if oldQ.UserCreated != nil && *oldQ.UserCreated != uid {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not authorized to duplicate this quotation")
		}
	}

	newID := repository.GenerateID()
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	ym := now.Format("200601")

	counter, err := h.repo.GetCounter(oldQ.PropertyID, "quotation", oldQ.QuotationType, ym)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate counter")
	}

	userInisialStr := "USR"
	if inisial, ok := userInisial.(string); ok && inisial != "" {
		userInisialStr = inisial
	}
	newQuotationID := repository.GenerateQuotationID(oldQ.QuotationType, counter, month, year, userInisialStr)

	err = h.repo.CreateRevision(oldID, newID, newQuotationID, now, uid, userInisialStr, "[DUPLICATE] ")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to duplicate: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, fiber.Map{
		"message":      "Quotation duplicated successfully",
		"new_id":       newID,
		"quotation_id": newQuotationID,
	})
}

func (h *QuotationHandler) SetDefault(c *fiber.Ctx) error {
	id := c.Params("id")
	revIDStr := c.Params("rev_id")
	if id == "" || revIDStr == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID and Revision ID are required")
	}

	revID, _ := strconv.Atoi(revIDStr)
	err := h.repo.SetDefault(id, revID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to set default revision: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Default revision updated successfully"})
}

func (h *QuotationHandler) ExportExcel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	search := c.Query("search")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	qType := c.Query("quotation_type")
	status := c.Query("status")
	userCreated := c.Query("sales_id")
	if userCreated == "" {
		userCreated = c.Query("user_created")
	}
	progress := c.Query("progress")

	if !h.userCanViewAllQuotations(userID) {
		userCreated = strconv.FormatUint(uint64(userID), 10)
	}

	page := 1
	limit := 100000 // Get all
	sortBy := c.Query("sort")
	sortDir := c.Query("order")

	items, _, err := h.repo.FindAll(search, fromDate, toDate, qType, status, userCreated, progress, page, limit, sortBy, sortDir)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve quotations for export: "+err.Error())
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetName := "Quotations"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"quotation_type", "customer_name", "quotation_id", "quotation_date", "subject",
		"total", "tax", "tax_value", "pph", "pph_value", "disc", "grand_total",
		"hpp_total", "profit", "profit_value", "valid_until", "commision",
		"notes", "total_grand", "status_name", "progress_name", "sales person",
	}

	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheetName, colName+"1", header)
	}

	for i, item := range items {
		row := i + 2
		qTypeStr := ""
		switch item.QuotationType {
		case 1:
			qTypeStr = "Project"
		case 2:
			qTypeStr = "Retail"
		case 3:
			qTypeStr = "Maintenance"
		case 4:
			qTypeStr = "Others"
		}

		custName := item.Customer.Name

		qDate := ""
		if item.QuotationDate != nil {
			qDate = item.QuotationDate.Format("2006-01-02")
		}

		validUntil := ""
		if item.ValidUntil != nil {
			validUntil = item.ValidUntil.Format("2006-01-02")
		}

		subject := ""
		if item.Subject != nil {
			subject = *item.Subject
		}

		total := 0.0
		if item.Total != nil {
			total = *item.Total
		}

		taxValue := item.TaxValue

		pphValue := item.PPhValue

		disc := 0.0
		if item.Disc != nil {
			disc = *item.Disc
		}

		grandTotal := 0.0
		if item.GrandTotal != nil {
			grandTotal = *item.GrandTotal
		}

		hppTotal := 0.0
		if item.HppTotal != nil {
			hppTotal = *item.HppTotal
		}

		profit := 0.0
		if item.Profit != nil {
			profit = *item.Profit
		}

		profitValue := 0.0
		if item.ProfitValue != nil {
			profitValue = *item.ProfitValue
		}

		notes := ""
		if item.Notes != nil {
			notes = *item.Notes
		}

		statusName := ""
		if item.StatusInfo != nil {
			statusName = item.StatusInfo.Name
		}

		progressName := ""
		if item.ProgressInfo != nil {
			progressName = item.ProgressInfo.Name
		}

		salesPerson := ""
		if item.SalesPerson != nil {
			salesPerson = item.SalesPerson.Name
		}

		totalGrandFormatted := fmt.Sprintf("Rp%.0f", grandTotal)

		values := []interface{}{
			qTypeStr, custName, item.QuotationID, qDate, subject,
			total, item.Tax, taxValue, item.PPh, pphValue, disc, grandTotal,
			hppTotal, profit, profitValue, validUntil, item.Commision,
			notes, totalGrandFormatted, statusName, progressName, salesPerson,
		}

		for j, val := range values {
			colName, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, row), val)
		}
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=quotations.xlsx")

	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate Excel file")
	}

	return nil
}

type TopCustomerResult struct {
	CustomerID    uint    `json:"customer_id"`
	CustomerName  string  `json:"customer_name"`
	GrandTotal    float64 `json:"grand_total"`
	Profit        float64 `json:"profit"`
	ProfitPercent float64 `json:"profit_percent"`
}

func (h *QuotationHandler) topCustomersByType(c *fiber.Ctx, fromDate, toDate string, quotationType int) ([]TopCustomerResult, error) {
	var results []TopCustomerResult
	err := h.repo.GetDB().Raw(`
		SELECT 
			c.id AS customer_id,
			c.name AS customer_name,
			COALESCE(SUM(m.grand_total), 0) AS grand_total,
			COALESCE(SUM(m.profit_value), 0) AS profit,
			COALESCE(SUM(m.grand_total * COALESCE(m.profit, 0)) / NULLIF(SUM(m.grand_total), 0), 0) AS profit_percent
		FROM quotation q
		JOIN quotation_master m ON m.id COLLATE utf8mb4_general_ci = q.id AND m.default_quot = 1
		JOIN customer c ON c.id = q.customer_id
		WHERE q.quotation_type = ?
			AND q.status = 3
			AND q.quotation_date BETWEEN ? AND ?
		GROUP BY c.id, c.name
		ORDER BY grand_total DESC
		LIMIT 10
	`, quotationType, fromDate, toDate).Scan(&results).Error

	return results, err
}

func (h *QuotationHandler) TopCustomers(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	project, err := h.topCustomersByType(c, fromDate, toDate, 1)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve project top customers: "+err.Error())
	}

	retail, err := h.topCustomersByType(c, fromDate, toDate, 2)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve retail top customers: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"retail":  retail,
		"project": project,
	})
}

func (h *QuotationHandler) TopCustomersExport(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	project, err := h.topCustomersByType(c, fromDate, toDate, 1)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve project top customers: "+err.Error())
	}

	retail, err := h.topCustomersByType(c, fromDate, toDate, 2)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve retail top customers: "+err.Error())
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	// Sheet 1: Retail
	retailSheet := "Retail"
	f.SetSheetName("Sheet1", retailSheet)

	retailHeaders := []string{
		"Rank", "Customer Name", "Grand Total", "Profit", "Margin %",
	}
	for i, header := range retailHeaders {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(retailSheet, colName+"1", header)
	}

	for i, item := range retail {
		row := i + 2
		values := []interface{}{
			i + 1,
			item.CustomerName,
			item.GrandTotal,
			item.Profit,
			fmt.Sprintf("%.2f%%", item.ProfitPercent),
		}
		for j, val := range values {
			colName, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(retailSheet, fmt.Sprintf("%s%d", colName, row), val)
		}
	}

	// Sheet 2: Project
	projectSheet := "Project"
	idx, err := f.NewSheet(projectSheet)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create sheet")
	}

	projectHeaders := []string{
		"Rank", "Customer Name", "Grand Total", "Profit", "Margin %",
	}
	for i, header := range projectHeaders {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(projectSheet, colName+"1", header)
	}

	for i, item := range project {
		row := i + 2
		values := []interface{}{
			i + 1,
			item.CustomerName,
			item.GrandTotal,
			item.Profit,
			fmt.Sprintf("%.2f%%", item.ProfitPercent),
		}
		for j, val := range values {
			colName, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(projectSheet, fmt.Sprintf("%s%d", colName, row), val)
		}
	}

	f.SetActiveSheet(idx)

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=top_customers.xlsx")

	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate Excel file")
	}

	return nil
}

type QuotationReportItem struct {
	QuotationID   string   `json:"quotation_id"`
	QuotationDate *string  `json:"quotation_date"`
	CustomerName  string   `json:"customer_name"`
	Subject       *string  `json:"subject"`
	ProgressName  *string  `json:"progress_name"`
	GrandTotal    float64  `json:"grand_total"`
	Profit        *float64 `json:"profit"`
	ProfitValue   *float64 `json:"profit_value"`
}

func (h *QuotationHandler) buildReportQuery(fromDate, toDate string, qType *int, search, status, userCreated, progress string, page, limit int) (*gorm.DB, int64) {
	query := h.repo.GetDB().Table("quotation q").
		Select(`q.quotation_id, q.quotation_date, c.name AS customer_name,
			q.subject, qp.name AS progress_name,
			q.grand_total, qm.profit, qm.profit_value`).
		Joins("JOIN customer c ON c.id = q.customer_id").
		Joins("LEFT JOIN quotation_progress qp ON qp.id = q.progress").
		Joins("LEFT JOIN quotation_master qm ON qm.id COLLATE utf8mb4_general_ci = q.id AND qm.default_quot = 1").
		Where("q.quotation_date BETWEEN ? AND ?", fromDate, toDate)

	if qType != nil {
		query = query.Where("q.quotation_type = ?", *qType)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(q.quotation_id LIKE ? OR q.subject LIKE ? OR c.name LIKE ?)", searchTerm, searchTerm, searchTerm)
	}

	if status != "" && status != "all" {
		if val, err := strconv.Atoi(status); err == nil {
			query = query.Where("q.status = ?", val)
		}
	}

	if userCreated != "" && userCreated != "all" {
		if val, err := strconv.Atoi(userCreated); err == nil {
			query = query.Where("q.sales_id = ?", val)
		}
	}

	if progress != "" && progress != "all" {
		if val, err := strconv.Atoi(progress); err == nil {
			query = query.Where("q.progress = ?", val)
		}
	}

	var total int64
	query.Count(&total)

	ordered := query.Order("q.quotation_date DESC, q.quotation_id DESC")

	if limit > 0 {
		ordered = ordered.Offset((page - 1) * limit).Limit(limit)
	}

	return ordered, total
}

func (h *QuotationHandler) QuotationsReport(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	status := c.Query("status")
	userCreated := c.Query("sales_id")
	if userCreated == "" {
		userCreated = c.Query("user_created")
	}
	progress := c.Query("progress")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	projectType := 1
	retailType := 2

	var all, retail, project []QuotationReportItem
	var allTotal, retailTotal, projectTotal int64

	allQ, allTotal := h.buildReportQuery(fromDate, toDate, nil, search, status, userCreated, progress, page, limit)
	if err := allQ.Find(&all).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve all quotations: "+err.Error())
	}

	retailQ, retailTotal := h.buildReportQuery(fromDate, toDate, &retailType, search, status, userCreated, progress, page, limit)
	if err := retailQ.Find(&retail).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve retail quotations: "+err.Error())
	}

	projectQ, projectTotal := h.buildReportQuery(fromDate, toDate, &projectType, search, status, userCreated, progress, page, limit)
	if err := projectQ.Find(&project).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve project quotations: "+err.Error())
	}

	users := h.fetchReportSalesPersons(c.Locals("user_id").(uint))

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"all":           all,
		"retail":        retail,
		"project":       project,
		"all_total":     allTotal,
		"retail_total":  retailTotal,
		"project_total": projectTotal,
		"page":          page,
		"limit":         limit,
		"users":         users,
	})
}

func (h *QuotationHandler) ReportSection(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	status := c.Query("status")
	userCreated := c.Query("sales_id")
	if userCreated == "" {
		userCreated = c.Query("user_created")
	}
	progress := c.Query("progress")
	section := c.Query("section")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var qType *int
	switch section {
	case "project":
		t := 1
		qType = &t
	case "retail":
		t := 2
		qType = &t
	}

	var items []QuotationReportItem
	q, total := h.buildReportQuery(fromDate, toDate, qType, search, status, userCreated, progress, page, limit)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve report section: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *QuotationHandler) writeReportSheet(f *excelize.File, sheet string, items []QuotationReportItem) {
	headers := []string{
		"quotation_id", "quotation_date", "customer_name",
		"subject", "progress_name", "grand_total", "profit", "profit_value",
	}
	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheet, colName+"1", header)
	}
	for i, item := range items {
		row := i + 2
		qDate := ""
		if item.QuotationDate != nil {
			qDate = *item.QuotationDate
		}
		subject := ""
		if item.Subject != nil {
			subject = *item.Subject
		}
		progressName := ""
		if item.ProgressName != nil {
			progressName = *item.ProgressName
		}
		profit := float64(0)
		if item.Profit != nil {
			profit = *item.Profit
		}
		profitValue := float64(0)
		if item.ProfitValue != nil {
			profitValue = *item.ProfitValue
		}
		values := []interface{}{
			item.QuotationID, qDate, item.CustomerName,
			subject, progressName, item.GrandTotal, profit, profitValue,
		}
		for j, val := range values {
			colName, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(sheet, fmt.Sprintf("%s%d", colName, row), val)
		}
	}
}

func (h *QuotationHandler) ReportExportExcel(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	status := c.Query("status")
	userCreated := c.Query("sales_id")
	if userCreated == "" {
		userCreated = c.Query("user_created")
	}
	progress := c.Query("progress")

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var all, retail, project []QuotationReportItem
	projectType := 1
	retailType := 2

	allQ, _ := h.buildReportQuery(fromDate, toDate, nil, search, status, userCreated, progress, 1, 0)
	if err := allQ.Find(&all).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve report data: "+err.Error())
	}
	retailQ, _ := h.buildReportQuery(fromDate, toDate, &retailType, search, status, userCreated, progress, 1, 0)
	if err := retailQ.Find(&retail).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve retail data: "+err.Error())
	}
	projectQ, _ := h.buildReportQuery(fromDate, toDate, &projectType, search, status, userCreated, progress, 1, 0)
	if err := projectQ.Find(&project).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve project data: "+err.Error())
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	f.SetSheetName("Sheet1", "All Quotation")
	f.NewSheet("Project Quotation")
	f.NewSheet("Retail Quotation")

	h.writeReportSheet(f, "All Quotation", all)
	h.writeReportSheet(f, "Project Quotation", project)
	h.writeReportSheet(f, "Retail Quotation", retail)

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=quotations-report.xlsx")

	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate Excel file")
	}

	return nil
}

type SalesSummaryItem struct {
	QuotationDate   *time.Time `json:"quotation_date"`
	QuotationTypeID int        `json:"quotation_type_id"`
	QuotationType   string     `json:"quotation_type"`
	CustomerName    string     `json:"customer_name"`
	GrandTotal      float64    `json:"grand_total"`
	ProfitValue     float64    `json:"profit_value"`
	MarginPercent   float64    `json:"margin_percent"`
}

func (h *QuotationHandler) buildSalesSummaryQuery(fromDate, toDate, search, quotationType string, page, limit int) (*gorm.DB, int64) {
	query := h.repo.GetDB().Table("quotation q").
		Select(`MAX(q.quotation_date) AS quotation_date,
			MAX(q.quotation_type) AS quotation_type_id,
			CASE MAX(q.quotation_type)
				WHEN 1 THEN 'Project'
				WHEN 2 THEN 'Retail'
				WHEN 3 THEN 'Maintenance'
				WHEN 4 THEN 'Others'
				ELSE 'Unknown'
			END AS quotation_type,
			c.id AS customer_id, c.name AS customer_name,
			COALESCE(SUM(q.grand_total),0) AS grand_total,
			COALESCE(SUM(q.profit_value),0) AS profit_value,
			CASE WHEN COALESCE(SUM(q.grand_total),0) = 0 THEN 0
				ELSE ROUND(COALESCE(SUM(q.profit_value),0) / COALESCE(SUM(q.grand_total),0) * 100, 2)
			END AS margin_percent`).
		Joins("JOIN customer c ON c.id = q.customer_id").
		Where("q.status = 3").
		Where("q.quotation_date BETWEEN ? AND ?", fromDate, toDate).
		Group("c.id, c.name")

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(c.name LIKE ?)", searchTerm)
	}

	if quotationType != "" && quotationType != "all" {
		if val, err := strconv.Atoi(quotationType); err == nil {
			query = query.Where("q.quotation_type = ?", val)
		}
	}

	var total int64
	query.Count(&total)

	ordered := query.Order("quotation_type_id ASC, grand_total DESC, c.name ASC")

	if limit > 0 {
		ordered = ordered.Offset((page - 1) * limit).Limit(limit)
	}

	return ordered, total
}

func (h *QuotationHandler) SalesSummaryByCustomer(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	quotationType := c.Query("quotation_type")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesSummaryItem
	q, total := h.buildSalesSummaryQuery(fromDate, toDate, search, quotationType, page, limit)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales summary: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *QuotationHandler) SalesSummaryByCustomerExport(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	quotationType := c.Query("quotation_type")

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesSummaryItem
	q, _ := h.buildSalesSummaryQuery(fromDate, toDate, search, quotationType, 1, 0)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales summary: "+err.Error())
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetName := "Sales Summary"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"quotation_date", "quotation_type", "customer_name",
		"grand_total", "profit_value", "margin_percent",
	}

	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheetName, colName+"1", header)
	}

	for i, item := range items {
		row := i + 2

		qDate := ""
		if item.QuotationDate != nil {
			qDate = item.QuotationDate.Format("2006-01-02")
		}

		values := []interface{}{
			qDate, item.QuotationType, item.CustomerName,
			item.GrandTotal, item.ProfitValue, item.MarginPercent,
		}

		for j, val := range values {
			colName, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, row), val)
		}
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=sales-summary-by-customer.xlsx")

	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate Excel file")
	}

	return nil
}

type SalesSummaryBySalesPersonItem struct {
	QuotationDate   *time.Time `json:"quotation_date"`
	QuotationTypeID int        `json:"quotation_type_id"`
	QuotationType   string     `json:"quotation_type"`
	SalesPerson     string     `json:"sales_person"`
	GrandTotal      float64    `json:"grand_total"`
	ProfitValue     float64    `json:"profit_value"`
	MarginPercent   float64    `json:"margin_percent"`
}

func (h *QuotationHandler) buildSalesSummaryBySalesPersonQuery(fromDate, toDate, search, salesId, quotationType string, page, limit int) (*gorm.DB, int64) {
	query := h.repo.GetDB().Table("quotation q").
		Select(`MAX(q.quotation_date) AS quotation_date,
			MAX(q.quotation_type) AS quotation_type_id,
			CASE MAX(q.quotation_type)
				WHEN 1 THEN 'Project'
				WHEN 2 THEN 'Retail'
				WHEN 3 THEN 'Maintenance'
				WHEN 4 THEN 'Others'
				ELSE 'Unknown'
			END AS quotation_type,
			u.id AS sales_person_id, u.name AS sales_person,
			COALESCE(SUM(q.grand_total),0) AS grand_total,
			COALESCE(SUM(q.profit_value),0) AS profit_value,
			CASE WHEN COALESCE(SUM(q.grand_total),0) = 0 THEN 0
				ELSE ROUND(COALESCE(SUM(q.profit_value),0) / COALESCE(SUM(q.grand_total),0) * 100, 2)
			END AS margin_percent`).
		Joins("LEFT JOIN users u ON u.id = q.sales_id").
		Where("q.status = 3").
		Where("q.quotation_date BETWEEN ? AND ?", fromDate, toDate).
		Group("u.id, u.name")

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(u.name LIKE ?)", searchTerm)
	}

	if salesId != "" && salesId != "all" {
		if val, err := strconv.Atoi(salesId); err == nil {
			query = query.Where("q.sales_id = ?", val)
		}
	}

	if quotationType != "" && quotationType != "all" {
		if val, err := strconv.Atoi(quotationType); err == nil {
			query = query.Where("q.quotation_type = ?", val)
		}
	}

	var total int64
	query.Count(&total)

	ordered := query.Order("quotation_type_id ASC, grand_total DESC, u.name ASC")

	if limit > 0 {
		ordered = ordered.Offset((page - 1) * limit).Limit(limit)
	}

	return ordered, total
}

func (h *QuotationHandler) SalesSummaryBySalesPerson(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	salesId := c.Query("sales_id")
	if salesId == "" {
		salesId = c.Query("user_created")
	}
	quotationType := c.Query("quotation_type")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesSummaryBySalesPersonItem
	q, total := h.buildSalesSummaryBySalesPersonQuery(fromDate, toDate, search, salesId, quotationType, page, limit)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales summary: "+err.Error())
	}

	users := h.fetchReportSalesPersons(userID)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
		"users": users,
	})
}

func (h *QuotationHandler) SalesSummaryBySalesPersonExport(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	salesId := c.Query("sales_id")
	if salesId == "" {
		salesId = c.Query("user_created")
	}
	quotationType := c.Query("quotation_type")

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesSummaryBySalesPersonItem
	q, _ := h.buildSalesSummaryBySalesPersonQuery(fromDate, toDate, search, salesId, quotationType, 1, 0)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales summary: "+err.Error())
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetName := "Sales Summary"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"quotation_date", "quotation_type", "sales_person",
		"grand_total", "profit_value", "margin_percent",
	}

	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheetName, colName+"1", header)
	}

	for i, item := range items {
		row := i + 2

		qDate := ""
		if item.QuotationDate != nil {
			qDate = item.QuotationDate.Format("2006-01-02")
		}

		values := []interface{}{
			qDate, item.QuotationType, item.SalesPerson,
			item.GrandTotal, item.ProfitValue, item.MarginPercent,
		}

		for j, val := range values {
			colName, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, row), val)
		}
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=sales-summary-by-sales-person.xlsx")

	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate Excel file")
	}

	return nil
}

type SalesDetailItem struct {
	QuotationID   string     `json:"quotation_id"`
	QuotationDate *time.Time `json:"quotation_date"`
	Subject       *string    `json:"subject"`
	QuotationType string     `json:"quotation_type"`
	SalesPerson   *string    `json:"sales_person"`
	CustomerName  *string    `json:"customer_name"`
	GrandTotal    float64    `json:"grand_total"`
	ProfitValue   float64    `json:"profit_value"`
	MarginPercent float64    `json:"margin_percent"`
}

type SalesDetailBySalesPersonItem struct {
	QuotationID   string     `json:"quotation_id"`
	QuotationDate *time.Time `json:"quotation_date"`
	Subject       *string    `json:"subject"`
	QuotationType string     `json:"quotation_type"`
	CustomerName  *string    `json:"customer_name"`
	GrandTotal    float64    `json:"grand_total"`
	ProfitValue   float64    `json:"profit_value"`
	MarginPercent float64    `json:"margin_percent"`
}

func (h *QuotationHandler) buildSalesDetailQuery(fromDate, toDate, search, salesId, customerId, quotationType string, page, limit int) (*gorm.DB, int64) {
	query := h.repo.GetDB().Table("quotation q").
		Select(`q.quotation_id, q.quotation_date, q.subject,
			CASE q.quotation_type
				WHEN 1 THEN 'Project'
				WHEN 2 THEN 'Retail'
				WHEN 3 THEN 'Maintenance'
				WHEN 4 THEN 'Others'
				ELSE 'Unknown'
			END AS quotation_type,
			u.name AS sales_person,
			c.name AS customer_name,
			COALESCE(q.grand_total,0) AS grand_total,
			COALESCE(q.profit_value,0) AS profit_value,
			CASE WHEN COALESCE(q.grand_total,0) = 0 THEN 0
				ELSE ROUND(COALESCE(q.profit_value,0) / COALESCE(q.grand_total,0) * 100, 2)
			END AS margin_percent`).
		Joins("LEFT JOIN users u ON u.id = q.sales_id").
		Joins("LEFT JOIN customer c ON c.id = q.customer_id").
		Where("q.status = 3").
		Where("q.quotation_date BETWEEN ? AND ?", fromDate, toDate)

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(q.quotation_id LIKE ? OR q.subject LIKE ?)", searchTerm, searchTerm)
	}

	if customerId != "" && customerId != "all" {
		if val, err := strconv.Atoi(customerId); err == nil {
			query = query.Where("q.customer_id = ?", val)
		}
	}

	if salesId != "" && salesId != "all" {
		if val, err := strconv.Atoi(salesId); err == nil {
			query = query.Where("q.sales_id = ?", val)
		}
	}

	if quotationType != "" && quotationType != "all" {
		if val, err := strconv.Atoi(quotationType); err == nil {
			query = query.Where("q.quotation_type = ?", val)
		}
	}

	var total int64
	query.Count(&total)

	ordered := query.Order("q.quotation_date ASC, q.quotation_id ASC")

	if limit > 0 {
		ordered = ordered.Offset((page - 1) * limit).Limit(limit)
	}

	return ordered, total
}

func (h *QuotationHandler) SalesDetailByCustomer(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	salesId := c.Query("sales_id")
	if salesId == "" {
		salesId = c.Query("user_created")
	}
	customerId := c.Query("customer_id")
	quotationType := c.Query("quotation_type")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesDetailItem
	q, total := h.buildSalesDetailQuery(fromDate, toDate, search, salesId, customerId, quotationType, page, limit)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales detail: "+err.Error())
	}

	users := h.fetchReportSalesPersons(userID)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
		"users": users,
	})
}

func (h *QuotationHandler) SalesDetailByCustomerExport(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	salesId := c.Query("sales_id")
	if salesId == "" {
		salesId = c.Query("user_created")
	}
	customerId := c.Query("customer_id")
	quotationType := c.Query("quotation_type")

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesDetailItem
	q, _ := h.buildSalesDetailQuery(fromDate, toDate, search, salesId, customerId, quotationType, 1, 0)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales detail: "+err.Error())
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetName := "Sales Detail"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"quotation_id", "quotation_date", "subject",
		"quotation_type", "customer_name",
		"grand_total", "profit_value", "margin_percent",
	}

	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheetName, colName+"1", header)
	}

	for i, item := range items {
		row := i + 2

		qDate := ""
		if item.QuotationDate != nil {
			qDate = item.QuotationDate.Format("2006-01-02")
		}

		subject := ""
		if item.Subject != nil {
			subject = *item.Subject
		}

		customerName := ""
		if item.CustomerName != nil {
			customerName = *item.CustomerName
		}

		values := []interface{}{
			item.QuotationID, qDate, subject,
			item.QuotationType, customerName,
			item.GrandTotal, item.ProfitValue, item.MarginPercent,
		}

		for j, val := range values {
			colName, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, row), val)
		}
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=sales-detail-by-customer.xlsx")

	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate Excel file")
	}

	return nil
}

func (h *QuotationHandler) buildSalesDetailBySalesPersonQuery(fromDate, toDate, search, salesId, quotationType string, page, limit int) (*gorm.DB, int64) {
	query := h.repo.GetDB().Table("quotation q").
		Select(`q.quotation_id, q.quotation_date, q.subject,
			CASE q.quotation_type
				WHEN 1 THEN 'Project'
				WHEN 2 THEN 'Retail'
				WHEN 3 THEN 'Maintenance'
				WHEN 4 THEN 'Others'
				ELSE 'Unknown'
			END AS quotation_type,
			c.name AS customer_name,
			COALESCE(q.grand_total,0) AS grand_total,
			COALESCE(q.profit_value,0) AS profit_value,
			CASE WHEN COALESCE(q.grand_total,0) = 0 THEN 0
				ELSE ROUND(COALESCE(q.profit_value,0) / COALESCE(q.grand_total,0) * 100, 2)
			END AS margin_percent`).
		Joins("LEFT JOIN customer c ON c.id = q.customer_id").
		Where("q.status = 3").
		Where("q.quotation_date BETWEEN ? AND ?", fromDate, toDate)

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(q.quotation_id LIKE ? OR q.subject LIKE ? OR c.name LIKE ?)", searchTerm, searchTerm, searchTerm)
	}

	if salesId != "" && salesId != "all" {
		if val, err := strconv.Atoi(salesId); err == nil {
			query = query.Where("q.sales_id = ?", val)
		}
	}

	if quotationType != "" && quotationType != "all" {
		if val, err := strconv.Atoi(quotationType); err == nil {
			query = query.Where("q.quotation_type = ?", val)
		}
	}

	var total int64
	countQuery := h.repo.GetDB().Table("(?) AS sub", query)
	countQuery.Count(&total)

	ordered := query.Order("q.quotation_date ASC, q.quotation_id ASC")

	if limit > 0 {
		ordered = ordered.Offset((page - 1) * limit).Limit(limit)
	}

	return ordered, total
}

func (h *QuotationHandler) SalesDetailBySalesPerson(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	salesId := c.Query("sales_id")
	if salesId == "" {
		salesId = c.Query("user_created")
	}
	quotationType := c.Query("quotation_type")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesDetailBySalesPersonItem
	q, total := h.buildSalesDetailBySalesPersonQuery(fromDate, toDate, search, salesId, quotationType, page, limit)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales detail: "+err.Error())
	}

	users := h.fetchReportSalesPersons(userID)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
		"users": users,
	})
}

func (h *QuotationHandler) SalesDetailBySalesPersonExport(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	salesId := c.Query("sales_id")
	if salesId == "" {
		salesId = c.Query("user_created")
	}
	quotationType := c.Query("quotation_type")

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesDetailBySalesPersonItem
	q, _ := h.buildSalesDetailBySalesPersonQuery(fromDate, toDate, search, salesId, quotationType, 1, 0)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales detail: "+err.Error())
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetName := "Sales Detail"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"quotation_id", "quotation_date", "subject",
		"quotation_type", "customer_name",
		"grand_total", "profit_value", "margin_percent",
	}

	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheetName, colName+"1", header)
	}

	for i, item := range items {
		row := i + 2

		qDate := ""
		if item.QuotationDate != nil {
			qDate = item.QuotationDate.Format("2006-01-02")
		}

		subject := ""
		if item.Subject != nil {
			subject = *item.Subject
		}

		customerName := ""
		if item.CustomerName != nil {
			customerName = *item.CustomerName
		}

		values := []interface{}{
			item.QuotationID, qDate, subject,
			item.QuotationType, customerName,
			item.GrandTotal, item.ProfitValue, item.MarginPercent,
		}

		for j, val := range values {
			colName, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, row), val)
		}
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=sales-detail-by-sales-person.xlsx")

	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate Excel file")
	}

	return nil
}

type SalesItemByCustomerItem struct {
	QuotationID   *string `json:"quotation_id"`
	PartNo        *string `json:"part_no"`
	Descriptions  *string `json:"descriptions"`
	SalesPerson   *string `json:"sales_person"`
	TotalQty      float64 `json:"total_qty"`
	HppTotal      float64 `json:"hpp_total"`
	GrandTotal    float64 `json:"grand_total"`
	ProfitValue   float64 `json:"profit_value"`
	MarginPercent float64 `json:"margin_percent"`
}

func (h *QuotationHandler) buildSalesItemByCustomerQuery(fromDate, toDate, search, customerId string, page, limit int) (*gorm.DB, int64) {
	query := h.repo.GetDB().Table("quotation_detail qd").
		Select(`MAX(q.quotation_id) AS quotation_id,
			qd.part_no, qd.descriptions,
			MAX(u.name) AS sales_person,
			COALESCE(SUM(qd.qty),0) AS total_qty,
			COALESCE(SUM(qd.hpp_total),0) AS hpp_total,
			COALESCE(SUM(qd.total),0) AS grand_total,
			COALESCE(SUM(qd.total),0) - COALESCE(SUM(qd.hpp_total),0) AS profit_value,
			CASE WHEN COALESCE(SUM(qd.total),0) = 0 THEN 0
				ELSE ROUND((COALESCE(SUM(qd.total),0) - COALESCE(SUM(qd.hpp_total),0)) / COALESCE(SUM(qd.total),0) * 100, 2)
			END AS margin_percent`).
		Joins("JOIN quotation q ON q.id COLLATE utf8mb4_general_ci = qd.id").
		Joins("JOIN quotation_master qm ON qm.id COLLATE utf8mb4_general_ci = qd.id AND qm.rev_id = qd.rev_id AND qm.default_quot = true").
		Joins("LEFT JOIN users u ON u.id = q.sales_id").
		Where("q.status = 3").
		Where("q.quotation_date BETWEEN ? AND ?", fromDate, toDate).
		Group("qd.part_no, qd.descriptions")

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(COALESCE(qd.part_no,'') LIKE ? OR COALESCE(qd.descriptions,'') LIKE ?)", searchTerm, searchTerm)
	}

	if customerId != "" && customerId != "all" {
		if val, err := strconv.Atoi(customerId); err == nil {
			query = query.Where("q.customer_id = ?", val)
		}
	}

	var total int64
	// Count from subquery for grouped results
	countQuery := h.repo.GetDB().Table("(?) AS sub", query)
	countQuery.Count(&total)

	ordered := query.Order("COALESCE(qd.part_no,'') ASC, COALESCE(qd.descriptions,'') ASC")

	if limit > 0 {
		ordered = ordered.Offset((page - 1) * limit).Limit(limit)
	}

	return ordered, total
}

func (h *QuotationHandler) SalesItemByCustomer(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	customerId := c.Query("customer_id")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesItemByCustomerItem
	q, total := h.buildSalesItemByCustomerQuery(fromDate, toDate, search, customerId, page, limit)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales item data: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *QuotationHandler) SalesItemByCustomerExport(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	customerId := c.Query("customer_id")

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesItemByCustomerItem
	q, _ := h.buildSalesItemByCustomerQuery(fromDate, toDate, search, customerId, 1, 0)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales item data: "+err.Error())
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetName := "Sales Item"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"quotation_id", "part_no", "descriptions", "sales_person", "total_qty",
		"hpp_total", "grand_total", "profit_value", "margin_percent",
	}

	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheetName, colName+"1", header)
	}

	for i, item := range items {
		row := i + 2

		quotationID := ""
		if item.QuotationID != nil {
			quotationID = *item.QuotationID
		}
		partNo := ""
		if item.PartNo != nil {
			partNo = *item.PartNo
		}
		desc := ""
		if item.Descriptions != nil {
			desc = *item.Descriptions
		}
		salesPerson := ""
		if item.SalesPerson != nil {
			salesPerson = *item.SalesPerson
		}

		values := []interface{}{
			quotationID, partNo, desc, salesPerson, item.TotalQty,
			item.HppTotal, item.GrandTotal, item.ProfitValue, item.MarginPercent,
		}

		for j, val := range values {
			colName, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, row), val)
		}
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=sales-item-by-customer.xlsx")

	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate Excel file")
	}

	return nil
}

func (h *QuotationHandler) buildSalesItemBySalesPersonQuery(fromDate, toDate, search, salesId string, page, limit int) (*gorm.DB, int64) {
	query := h.repo.GetDB().Table("quotation_detail qd").
		Select(`MAX(q.quotation_id) AS quotation_id,
			qd.part_no, qd.descriptions,
			MAX(u.name) AS sales_person,
			COALESCE(SUM(qd.qty),0) AS total_qty,
			COALESCE(SUM(qd.hpp_total),0) AS hpp_total,
			COALESCE(SUM(qd.total),0) AS grand_total,
			COALESCE(SUM(qd.total),0) - COALESCE(SUM(qd.hpp_total),0) AS profit_value,
			CASE WHEN COALESCE(SUM(qd.total),0) = 0 THEN 0
				ELSE ROUND((COALESCE(SUM(qd.total),0) - COALESCE(SUM(qd.hpp_total),0)) / COALESCE(SUM(qd.total),0) * 100, 2)
			END AS margin_percent`).
		Joins("JOIN quotation q ON q.id COLLATE utf8mb4_general_ci = qd.id").
		Joins("JOIN quotation_master qm ON qm.id COLLATE utf8mb4_general_ci = qd.id AND qm.rev_id = qd.rev_id AND qm.default_quot = true").
		Joins("LEFT JOIN users u ON u.id = q.sales_id").
		Where("q.status = 3").
		Where("q.quotation_date BETWEEN ? AND ?", fromDate, toDate).
		Group("qd.part_no, qd.descriptions")

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(COALESCE(qd.part_no,'') LIKE ? OR COALESCE(qd.descriptions,'') LIKE ?)", searchTerm, searchTerm)
	}

	if salesId != "" && salesId != "all" {
		if val, err := strconv.Atoi(salesId); err == nil {
			query = query.Where("q.sales_id = ?", val)
		}
	}

	var total int64
	countQuery := h.repo.GetDB().Table("(?) AS sub", query)
	countQuery.Count(&total)

	ordered := query.Order("COALESCE(qd.part_no,'') ASC, COALESCE(qd.descriptions,'') ASC")

	if limit > 0 {
		ordered = ordered.Offset((page - 1) * limit).Limit(limit)
	}

	return ordered, total
}

func (h *QuotationHandler) SalesItemBySalesPerson(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	salesId := c.Query("sales_id")
	if salesId == "" {
		salesId = c.Query("user_created")
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesItemByCustomerItem
	q, total := h.buildSalesItemBySalesPersonQuery(fromDate, toDate, search, salesId, page, limit)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales item data: "+err.Error())
	}

	users := h.fetchReportSalesPersons(userID)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
		"users": users,
	})
}

func (h *QuotationHandler) SalesItemBySalesPersonExport(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	salesId := c.Query("sales_id")
	if salesId == "" {
		salesId = c.Query("user_created")
	}

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []SalesItemByCustomerItem
	q, _ := h.buildSalesItemBySalesPersonQuery(fromDate, toDate, search, salesId, 1, 0)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales item data: "+err.Error())
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetName := "Sales Item"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"quotation_id", "part_no", "descriptions", "sales_person", "total_qty",
		"hpp_total", "grand_total", "profit_value", "margin_percent",
	}

	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheetName, colName+"1", header)
	}

	for i, item := range items {
		row := i + 2

		quotationID := ""
		if item.QuotationID != nil {
			quotationID = *item.QuotationID
		}
		partNo := ""
		if item.PartNo != nil {
			partNo = *item.PartNo
		}
		desc := ""
		if item.Descriptions != nil {
			desc = *item.Descriptions
		}
		salesPerson := ""
		if item.SalesPerson != nil {
			salesPerson = *item.SalesPerson
		}

		values := []interface{}{
			quotationID, partNo, desc, salesPerson, item.TotalQty,
			item.HppTotal, item.GrandTotal, item.ProfitValue, item.MarginPercent,
		}

		for j, val := range values {
			colName, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, row), val)
		}
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=sales-item-by-sales-person.xlsx")

	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate Excel file")
	}

	return nil
}

type SalesChartMonthlyItem struct {
	Month      string  `json:"month"`
	GrandTotal float64 `json:"grand_total"`
}

type SalesChartMonthlySeriesItem struct {
	Month      string  `json:"month"`
	GrandTotal float64 `json:"grand_total"`
	StatusName string  `json:"status_name"`
}

type SalesChartItemDetail struct {
	PartNo     *string `json:"part_no"`
	TotalQty   float64 `json:"total_qty"`
	TotalSales float64 `json:"total_sales"`
}

func (h *QuotationHandler) SalesChartsByCustomer(c *fiber.Ctx) error {
	customerID := c.Query("customer_id")
	year := c.Query("year")

	if customerID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "customer_id is required")
	}
	if year == "" {
		year = strconv.Itoa(time.Now().Year())
	}

	var monthly []SalesChartMonthlyItem
	if err := h.repo.GetDB().Table("quotation").
		Select(`DATE_FORMAT(quotation_date, '%Y-%m') AS month, COALESCE(SUM(grand_total),0) AS grand_total`).
		Where("customer_id = ?", customerID).
		Where("status = 3").
		Where("YEAR(quotation_date) = ?", year).
		Group("month").
		Order("month ASC").
		Find(&monthly).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve chart data: "+err.Error())
	}

	var items []SalesChartItemDetail
	if err := h.repo.GetDB().Table("quotation_detail qd").
		Select(`qd.part_no, COALESCE(SUM(qd.qty),0) AS total_qty, COALESCE(SUM(qd.total),0) AS total_sales`).
		Joins("JOIN quotation q ON q.id COLLATE utf8mb4_general_ci = qd.id").
		Joins("JOIN quotation_master qm ON qm.id COLLATE utf8mb4_general_ci = qd.id AND qm.rev_id = qd.rev_id AND qm.default_quot = true").
		Where("q.customer_id = ?", customerID).
		Where("q.status = 3").
		Where("YEAR(q.quotation_date) = ?", year).
		Group("qd.part_no, qd.descriptions").
		Order("total_sales DESC").
		Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve item details: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"monthly": monthly,
		"items":   items,
	})
}

type SalesChartsBySalesPersonItem struct {
	CustomerName  string  `json:"customer_name"`
	GrandTotal    float64 `json:"grand_total"`
	HppTotal      float64 `json:"hpp_total"`
	ProfitValue   float64 `json:"profit_value"`
	MarginPercent float64 `json:"margin_percent"`
}

func (h *QuotationHandler) SalesChartsBySalesPerson(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	salesID := c.Query("sales_id")
	if salesID == "" {
		salesID = c.Query("user_created")
	}
	year := c.Query("year")

	users := h.fetchReportSalesPersons(userID)

	if salesID == "" {
		return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
			"monthly":           []SalesChartMonthlySeriesItem{},
			"items":             []SalesChartsBySalesPersonItem{},
			"on_progress_items": []SalesChartsBySalesPersonItem{},
			"users":             users,
		})
	}
	if year == "" {
		year = strconv.Itoa(time.Now().Year())
	}

	type statusFilter struct {
		Name string
	}
	statuses := []statusFilter{
		{Name: "All Quotation"},
		{Name: "On Progress"},
		{Name: "P.O."},
		{Name: "Cancel"},
		{Name: "Decline"},
	}
	// Map status names to IDs
	statusIDs := map[string]int{
		"On Progress": 1,
		"Decline":     2,
		"P.O.":        3,
		"Cancel":      4,
	}

	var monthly []SalesChartMonthlySeriesItem
	for _, sf := range statuses {
		q := h.repo.GetDB().Table("quotation").
			Select(`DATE_FORMAT(quotation_date, '%Y-%m') AS month, COALESCE(SUM(grand_total),0) AS grand_total`).
			Where("user_created = ?", salesID).
			Where("YEAR(quotation_date) = ?", year).
			Group("month").
			Order("month ASC")
		if sf.Name == "All Quotation" {
			// All statuses
		} else if id, ok := statusIDs[sf.Name]; ok {
			q = q.Where("status = ?", id)
		}
		var rows []SalesChartMonthlySeriesItem
		if err := q.Find(&rows).Error; err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve chart data: "+err.Error())
		}
		for i := range rows {
			rows[i].StatusName = sf.Name
		}
		monthly = append(monthly, rows...)
	}

	var items []SalesChartsBySalesPersonItem
	if err := h.repo.GetDB().Table("quotation q").
		Select(`c.name AS customer_name, COALESCE(SUM(q.grand_total),0) AS grand_total, COALESCE(SUM(q.hpp_total),0) AS hpp_total, COALESCE(SUM(q.grand_total - q.hpp_total),0) AS profit_value, CASE WHEN SUM(q.grand_total) > 0 THEN ROUND((SUM(q.grand_total - q.hpp_total) / SUM(q.grand_total)) * 100, 2) ELSE 0 END AS margin_percent`).
		Joins("JOIN customer c ON c.id = q.customer_id").
		Where("q.user_created = ?", salesID).
		Where("q.status = 3").
		Where("YEAR(q.quotation_date) = ?", year).
		Group("c.id, c.name").
		Order("grand_total DESC").
		Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve item details: "+err.Error())
	}

	var onProgressItems []SalesChartsBySalesPersonItem
	if err := h.repo.GetDB().Table("quotation q").
		Select(`c.name AS customer_name, COALESCE(SUM(q.grand_total),0) AS grand_total, COALESCE(SUM(q.hpp_total),0) AS hpp_total, COALESCE(SUM(q.grand_total - q.hpp_total),0) AS profit_value, CASE WHEN SUM(q.grand_total) > 0 THEN ROUND((SUM(q.grand_total - q.hpp_total) / SUM(q.grand_total)) * 100, 2) ELSE 0 END AS margin_percent`).
		Joins("JOIN customer c ON c.id = q.customer_id").
		Where("q.user_created = ?", salesID).
		Where("q.status = 1").
		Where("YEAR(q.quotation_date) = ?", year).
		Group("c.id, c.name").
		Order("grand_total DESC").
		Find(&onProgressItems).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve on progress details: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"monthly":           monthly,
		"items":             items,
		"on_progress_items": onProgressItems,
		"users":             users,
	})
}

// ChartProgressByYear returns monthly grand total aggregated by quotation status name for a given year
type ChartProgressItem struct {
	Month      string  `json:"month"`
	StatusName string  `json:"status_name"`
	GrandTotal float64 `json:"grand_total"`
	Count      int64   `json:"count"`
}

func (h *QuotationHandler) ChartProgressByYear(c *fiber.Ctx) error {
	year := c.Query("year")
	if year == "" {
		year = strconv.Itoa(time.Now().Year())
	}

	var items []ChartProgressItem
	err := h.repo.GetDB().Table("quotation q").
		Select(`DATE_FORMAT(q.quotation_date, '%Y-%m') AS month,
			qs.name AS status_name,
			COALESCE(SUM(q.grand_total), 0) AS grand_total,
			COUNT(*) AS count`).
		Joins("JOIN quotation_status qs ON qs.id = q.status").
		Where("YEAR(q.quotation_date) = ?", year).
		Where("q.quotation_date IS NOT NULL").
		Group("month, qs.name").
		Order("month ASC, status_name ASC").
		Find(&items).Error

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve chart data: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"year":  year,
	})
}

// ChartPOByTypeItem represents PO grand total grouped by quotation type
type ChartPOByTypeItem struct {
	TypeName   string  `json:"type_name"`
	GrandTotal float64 `json:"grand_total"`
}

// ChartPOBySalesItem represents PO grand total grouped by sales person
type ChartPOBySalesItem struct {
	SalesName  string  `json:"sales_name"`
	GrandTotal float64 `json:"grand_total"`
}

// ChartPOAnalysis returns PO quotation data grouped by type and sales person for a given year
func (h *QuotationHandler) ChartPOAnalysis(c *fiber.Ctx) error {
	year := c.Query("year")
	if year == "" {
		year = strconv.Itoa(time.Now().Year())
	}

	// PO by Quotation Type
	var byType []ChartPOByTypeItem
	err := h.repo.GetDB().Table("quotation q").
		Select(`CASE 
				WHEN q.quotation_type = 1 THEN 'Project'
				WHEN q.quotation_type = 2 THEN 'Retail'
				WHEN q.quotation_type = 3 THEN 'Maintenance'
				WHEN q.quotation_type = 4 THEN 'Others'
				ELSE 'Unknown'
			END AS type_name,
			COALESCE(SUM(q.grand_total), 0) AS grand_total`).
		Where("q.status = 3").
		Where("YEAR(q.quotation_date) = ?", year).
		Where("q.quotation_date IS NOT NULL").
		Group("q.quotation_type").
		Order("grand_total DESC").
		Find(&byType).Error

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve PO by type data: "+err.Error())
	}

	// PO by Sales Person
	var bySales []ChartPOBySalesItem
	err = h.repo.GetDB().Table("quotation q").
		Select(`COALESCE(u.name, 'Unknown') AS sales_name,
			COALESCE(SUM(q.grand_total), 0) AS grand_total`).
		Joins("LEFT JOIN users u ON u.id = q.sales_id").
		Where("q.status = 3").
		Where("YEAR(q.quotation_date) = ?", year).
		Where("q.quotation_date IS NOT NULL").
		Group("q.sales_id, u.name").
		Order("grand_total DESC").
		Find(&bySales).Error

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve PO by sales data: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"by_type":  byType,
		"by_sales": bySales,
		"year":     year,
	})
}

// QuotationStats returns count and grand total by status for a given year
func (h *QuotationHandler) QuotationStats(c *fiber.Ctx) error {
	year := c.Query("year")
	if year == "" {
		year = strconv.Itoa(time.Now().Year())
	}

	type statusStat struct {
		StatusName string  `json:"status_name"`
		Count      int64   `json:"count"`
		GrandTotal float64 `json:"grand_total"`
	}

	var rows []statusStat
	err := h.repo.GetDB().Table("quotation q").
		Select(`COALESCE(qs.name, 'Unknown') AS status_name,
			COUNT(*) AS count,
			COALESCE(SUM(q.grand_total), 0) AS grand_total`).
		Joins("LEFT JOIN quotation_status qs ON qs.id = q.status").
		Where("YEAR(q.quotation_date) = ?", year).
		Where("q.quotation_date IS NOT NULL").
		Group("q.status, qs.name").
		Order("q.status ASC").
		Find(&rows).Error

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve quotation stats: "+err.Error())
	}

	// Compute overall totals by summing per-status rows
	var overallCount int64
	var overallTotal float64
	for _, r := range rows {
		overallCount += r.Count
		overallTotal += r.GrandTotal
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items":    rows,
		"total":    overallCount,
		"total_gt": overallTotal,
		"year":     year,
	})
}

// NeedFollowup returns quotations needing follow-up for the logged-in user
func (h *QuotationHandler) NeedFollowup(c *fiber.Ctx) error {
	uid := getUID(c)
	if uid == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	now := time.Now()

	type FollowupItem struct {
		ID            string     `json:"id"`
		QuotationID   string     `json:"quotation_id"`
		QuotationDate *time.Time `json:"quotation_date"`
		CustomerName  string     `json:"customer_name"`
		Subject       *string    `json:"subject"`
		GrandTotal    *float64   `json:"grand_total"`
		NextFollowup  *time.Time `json:"next_followup"`
	}

	var items []FollowupItem
	err := h.repo.GetDB().Table("quotation q").
		Select(`q.id, q.quotation_id, q.quotation_date,
			COALESCE(c.name, '') AS customer_name,
			q.subject, q.grand_total, q.next_followup`).
		Joins("LEFT JOIN customer c ON c.id = q.customer_id").
		Where("q.next_followup IS NOT NULL").
		Where("q.status = ?", 1).
		Where("(q.sales_id = ? OR q.user_created = ?)", uid, uid).
		Where("q.next_followup BETWEEN ? AND ?", threeMonthsAgo, now).
		Order("q.next_followup ASC").
		Limit(20).
		Find(&items).Error

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve follow-up data: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
	})
}
