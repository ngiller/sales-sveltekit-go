package handlers

import (
	"backend/internal/repository"
	"backend/internal/utils"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
	"github.com/gofiber/fiber/v2"
)

type StockHandler struct {
	repo *repository.StockRepository
}

func NewStockHandler(repo *repository.StockRepository) *StockHandler {
	return &StockHandler{repo: repo}
}

func (h *StockHandler) GetAllProducts(c *fiber.Ctx) error {
	search := c.Query("search", "")
	categoryID := c.Query("category_id", "")
	brandID := c.Query("brand_id", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	sortBy := c.Query("sort", "product_name")
	sortDir := c.Query("order", "asc")

	items, total, err := h.repo.FindAllProducts(search, categoryID, brandID, page, limit, sortBy, sortDir)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve live stocks: "+err.Error())
	}

	// Calculate total_cogs
	totalCogs := 0.0
	for _, item := range items {
		if endStock, ok := item["end_stock"].(int64); ok {
			if lastSellingPrice, ok := item["last_selling_price"].(string); ok {
				// Parse price string to float
				if price, err := strconv.ParseFloat(lastSellingPrice, 64); err == nil {
					totalCogs += float64(endStock) * price
				}
			}
		}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items":      items,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_cogs": totalCogs,
	})
}

func (h *StockHandler) ExportToExcel(c *fiber.Ctx) error {
	search := c.Query("search", "")
	categoryID := c.Query("category_id", "")
	brandID := c.Query("brand_id", "")

	// Get all items without pagination for export
	items, _, err := h.repo.FindAllProducts(search, categoryID, brandID, 1, 100000, "product_name", "asc")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve data for export: "+err.Error())
	}

	// Create Excel file
	f := excelize.NewFile()
	// Rename default sheet to "Live Stocks"
	f.SetSheetName("Sheet1", "Live Stocks")
	sheetName := "Live Stocks"
	
	// Set headers
	headers := []string{"No", "Product Code", "Category", "Brand", "Product Name", "Stock", "Last Buy Date", "Buy Price", "Selling Price"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Set data
	for rowIdx, item := range items {
		rowNum := rowIdx + 2
		
		// No (auto number)
		cell, _ := excelize.CoordinatesToCellName(1, rowNum)
		f.SetCellValue(sheetName, cell, rowIdx+1)
		
		// Product Code
		cell, _ = excelize.CoordinatesToCellName(2, rowNum)
		f.SetCellValue(sheetName, cell, getStringValue(item, "product_code"))
		
		// Category
		cell, _ = excelize.CoordinatesToCellName(3, rowNum)
		f.SetCellValue(sheetName, cell, getStringValue(item, "product_category_name"))
		
		// Brand
		cell, _ = excelize.CoordinatesToCellName(4, rowNum)
		f.SetCellValue(sheetName, cell, getStringValue(item, "brand_name"))
		
		// Product Name
		cell, _ = excelize.CoordinatesToCellName(5, rowNum)
		f.SetCellValue(sheetName, cell, getStringValue(item, "product_name"))
		
		// Stock
		cell, _ = excelize.CoordinatesToCellName(6, rowNum)
		f.SetCellValue(sheetName, cell, getIntValue(item, "end_stock"))
		
		// Last Buy Date
		cell, _ = excelize.CoordinatesToCellName(7, rowNum)
		f.SetCellValue(sheetName, cell, formatDate(getStringValue(item, "last_buy_date")))
		
		// Buy Price
		cell, _ = excelize.CoordinatesToCellName(8, rowNum)
		f.SetCellValue(sheetName, cell, getFloatValue(item, "last_buy_price"))
		
		// Selling Price
		cell, _ = excelize.CoordinatesToCellName(9, rowNum)
		f.SetCellValue(sheetName, cell, getFloatValue(item, "last_selling_price"))
	}

	// Auto fit column width
	for i := range headers {
		col := string(rune('A' + i))
		f.SetColWidth(sheetName, col, col, 15)
	}

	// Set content type and headers
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=live_stocks.xlsx")

	// Write to response
	file, err := f.WriteToBuffer()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate Excel: "+err.Error())
	}

	return c.Send(file.Bytes())
}

func getStringValue(item map[string]interface{}, key string) string {
	if val, ok := item[key]; ok && val != nil {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func getIntValue(item map[string]interface{}, key string) int {
	if val, ok := item[key]; ok && val != nil {
		switch v := val.(type) {
		case int64:
			return int(v)
		case int:
			return v
		case float64:
			return int(v)
		case string:
			i, _ := strconv.Atoi(v)
			return i
		}
	}
	return 0
}

func getFloatValue(item map[string]interface{}, key string) float64 {
	if val, ok := item[key]; ok && val != nil {
		switch v := val.(type) {
		case float64:
			return v
		case int64:
			return float64(v)
		case int:
			return float64(v)
		case string:
			f, _ := strconv.ParseFloat(v, 64)
			return f
		}
	}
	return 0
}

func formatDate(dateStr string) string {
	if len(dateStr) >= 10 {
		return dateStr[:10]
	}
	return dateStr
}
