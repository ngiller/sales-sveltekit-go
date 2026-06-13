package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type QuotationFollowupHandler struct {
	repo *repository.QuotationFollowupRepository
}

func NewQuotationFollowupHandler(repo *repository.QuotationFollowupRepository) *QuotationFollowupHandler {
	return &QuotationFollowupHandler{repo: repo}
}

func (h *QuotationFollowupHandler) FindAllByQuotation(c *fiber.Ctx) error {
	quotationID := c.Query("quotation_id")
	if quotationID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "quotation_id is required")
	}

	items, err := h.repo.FindAllByQuotationID(quotationID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve follow-ups: "+err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusOK, items)
}

func (h *QuotationFollowupHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	lineID, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid line_id format")
	}
	quotationID := c.Query("quotation_id")
	if quotationID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "quotation_id is required")
	}

	item, err := h.repo.FindByID(quotationID, lineID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Follow-up not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, item)
}

func (h *QuotationFollowupHandler) Create(c *fiber.Ctx) error {
	var item models.QuotationFollowup
	if err := c.BodyParser(&item); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	// Fetch property_id and default rev_id from Quotation
	var quot models.Quotation
	if err := h.repo.GetDB().Where("id = ?", item.ID).First(&quot).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid quotation ID (parent not found)")
	}
	item.PropertyID = quot.PropertyID

	var master models.QuotationMaster
	if err := h.repo.GetDB().Where("id = ? AND default_quot = ?", item.ID, true).First(&master).Error; err == nil {
		item.RevID = master.RevID
	} else {
		item.RevID = 1 // Fallback
	}

	// Set current user as followup_by if not set
	uid := c.Locals("user_id").(uint)
	if item.FollowupBy == nil {
		item.FollowupBy = &uid
	}

	// Auto-calculate the next sequential line_id for this quotation ID (item.ID)
	var maxLine int
	h.repo.GetDB().Model(&models.QuotationFollowup{}).
		Where("id = ?", item.ID).
		Select("COALESCE(MAX(line_id), 0)").
		Row().Scan(&maxLine)
	item.LineID = maxLine + 1

	if err := h.repo.Create(&item); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create follow-up: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, item)
}

func (h *QuotationFollowupHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	lineID, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid line_id format")
	}

	var req models.QuotationFollowup
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	item, err := h.repo.FindByID(req.ID, lineID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Follow-up not found")
	}

	item.Status = req.Status
	item.Progress = req.Progress
	item.FollowupDate = req.FollowupDate
	item.NextFollowup = req.NextFollowup
	item.Notes = req.Notes
	item.PoNo = req.PoNo
	item.PoDate = req.PoDate
	item.PoFile = req.PoFile

	if err := h.repo.Update(item); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update follow-up: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, item)
}

func (h *QuotationFollowupHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	lineID, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid line_id format")
	}
	quotationID := c.Query("quotation_id")
	if quotationID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "quotation_id is required")
	}

	if err := h.repo.Delete(quotationID, lineID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete follow-up")
	}

	// Check remaining follow-ups
	var count int64
	h.repo.GetDB().Model(&models.QuotationFollowup{}).Where("id = ?", quotationID).Count(&count)

	if count == 0 {
		// No follow-ups left → revert to "On Progress" (status=1) and "Lead" (progress=2)
		h.repo.GetDB().Model(&models.Quotation{}).Where("id = ?", quotationID).Updates(map[string]interface{}{
			"status":   1,
			"progress": 2,
		})
	} else {
		// Get the last follow-up (highest line_id) and sync its status/progress to quotation
		var last struct {
			Status   *uint `json:"status"`
			Progress *uint `json:"progress"`
		}
		h.repo.GetDB().Model(&models.QuotationFollowup{}).
			Select("status, progress").
			Where("id = ?", quotationID).
			Order("line_id DESC").
			Limit(1).
			Take(&last)
		updates := map[string]interface{}{}
		if last.Status != nil {
			updates["status"] = *last.Status
		}
		if last.Progress != nil {
			updates["progress"] = *last.Progress
		}
		if len(updates) > 0 {
			h.repo.GetDB().Model(&models.Quotation{}).Where("id = ?", quotationID).Updates(updates)
		}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Follow-up successfully deleted"})
}

func (h *QuotationFollowupHandler) UploadPO(c *fiber.Ctx) error {
	file, err := c.FormFile("po_file")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded")
	}

	if msg := utils.ValidateFile(file, utils.MergeExts(utils.AllowedImageExts(), utils.AllowedDocumentExts()), 10*1024*1024); msg != "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, msg)
	}

	// Create directory if not exists
	uploadDir := filepath.Join("uploads", "po_files")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create upload directory")
	}

	extension := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), extension)
	savePath := filepath.Join(uploadDir, newFilename)

	if err := c.SaveFile(file, savePath); err != nil {
		log.Printf("ERROR: Failed to save PO file: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save file")
	}

	fileURL := fmt.Sprintf("%s/%s", c.BaseURL(), filepath.ToSlash(savePath))
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message":   "PO file uploaded successfully",
		"url":       fileURL,
		"file_name": file.Filename,
	})
}

type FollowupReportItem struct {
	QuotationID   string     `json:"quotation_id"`
	QuotationDate *time.Time `json:"quotation_date"`
	FollowupDate  *time.Time `json:"followup_date"`
	NextFollowup  *time.Time `json:"next_followup"`
	CustomerName  string     `json:"customer_name"`
	Subject       *string    `json:"subject"`
	ProgressName  *string    `json:"progress_name"`
	GrandTotal    *float64   `json:"grand_total"`
	SalesPerson   *string    `json:"sales_person"`
}

func (h *QuotationFollowupHandler) buildFollowupReportQuery(fromDate, toDate string, search, salesId, progress, quotationType string, page, limit int) (*gorm.DB, int64) {
	// Query from quotation (all non-PO), data di-sync dari quotation_followup
	query := h.repo.GetDB().Table("quotation q").
		Select(`q.quotation_id, q.quotation_date, q.followup_date, q.next_followup,
			c.name AS customer_name, q.subject, qp.name AS progress_name,
			q.grand_total, u.name AS sales_person`).
		Joins("JOIN customer c ON c.id = q.customer_id").
		Joins("LEFT JOIN users u ON u.id = q.sales_id").
		Joins("LEFT JOIN quotation_progress qp ON qp.id = q.progress").
		Where("(q.po_no IS NULL OR q.po_no = '')").
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

	if progress != "" && progress != "all" {
		if val, err := strconv.Atoi(progress); err == nil {
			query = query.Where("q.progress = ?", val)
		}
	}

	if quotationType != "" && quotationType != "all" {
		if val, err := strconv.Atoi(quotationType); err == nil {
			query = query.Where("q.quotation_type = ?", val)
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

func (h *QuotationFollowupHandler) FollowupsReport(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	salesId := c.Query("sales_id")
	if salesId == "" {
		salesId = c.Query("user_created")
	}
	progress := c.Query("progress")
	quotationType := c.Query("quotation_type")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []FollowupReportItem
	q, total := h.buildFollowupReportQuery(fromDate, toDate, search, salesId, progress, quotationType, page, limit)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve followup report: "+err.Error())
	}

	var userItems []struct {
		ID      uint    `json:"id"`
		Name    string  `json:"name"`
		Inisial *string `json:"inisial"`
	}
	h.repo.GetDB().Model(&models.User{}).Select("id, name, inisial").Order("name ASC").Find(&userItems)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
		"users": userItems,
	})
}

func (h *QuotationFollowupHandler) FollowupsReportExport(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	search := c.Query("search")
	salesId := c.Query("sales_id")
	if salesId == "" {
		salesId = c.Query("user_created")
	}
	progress := c.Query("progress")
	quotationType := c.Query("quotation_type")

	if fromDate == "" || toDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "from_date and to_date are required")
	}

	var items []FollowupReportItem
	q, _ := h.buildFollowupReportQuery(fromDate, toDate, search, salesId, progress, quotationType, 1, 0)
	if err := q.Find(&items).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve followup report: "+err.Error())
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetName := "Followups Report"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"quotation_id", "quotation_date", "next_followup", "customer_name",
		"subject", "progress_name", "grand_total", "sales_person",
	}

	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheetName, colName+"1", header)
	}

	for i, item := range items {
		row := i + 2

		quotationDate := ""
		if item.QuotationDate != nil {
			quotationDate = item.QuotationDate.Format("2006-01-02")
		}

		nextFollowup := ""
		if item.NextFollowup != nil {
			nextFollowup = item.NextFollowup.Format("2006-01-02")
		}

		subject := ""
		if item.Subject != nil {
			subject = *item.Subject
		}

		progressName := ""
		if item.ProgressName != nil {
			progressName = *item.ProgressName
		}

		salesPerson := ""
		if item.SalesPerson != nil {
			salesPerson = *item.SalesPerson
		}

		grandTotal := 0.0
		if item.GrandTotal != nil {
			grandTotal = *item.GrandTotal
		}

		values := []interface{}{
			item.QuotationID, quotationDate, nextFollowup, item.CustomerName,
			subject, progressName, grandTotal, salesPerson,
		}

		for j, val := range values {
			colName, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, row), val)
		}
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=followups-report.xlsx")

	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate Excel file")
	}

	return nil
}
