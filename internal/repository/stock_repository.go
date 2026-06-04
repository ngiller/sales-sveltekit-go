package repository

import (
	"gorm.io/gorm"
)

type StockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) FindAllProducts(search string, categoryID, brandID string, page, limit int, sortBy, sortDir string) ([]map[string]interface{}, int64, error) {
	var items []map[string]interface{}
	var total int64

	// Count total products
	countQuery := r.db.Table("product")
	if search != "" {
		searchTerm := "%" + search + "%"
		countQuery = countQuery.Where("product_name LIKE ? OR product_code LIKE ?", searchTerm, searchTerm)
	}
	if categoryID != "" {
		countQuery = countQuery.Where("category_id = ?", categoryID)
	}
	if brandID != "" {
		countQuery = countQuery.Where("brand_id = ?", brandID)
	}
	countQuery.Count(&total)

	// Build raw SQL query with JOIN to get stock data
	// COALESCE handles NULL values
	sql := `
		SELECT 
			p.id,
			p.product_name,
			p.product_code,
			p.brand_id,
			p.category_id,
			b.brand_name,
			pc.product_category_name,
			COALESCE(ps.end_stock, 0) as end_stock,
			COALESCE(ps.last_buy_price, '0') as last_buy_price,
			COALESCE(ps.last_selling_price, '0') as last_selling_price,
			ps.last_buy_date
		FROM product p
		LEFT JOIN brand b ON b.id = p.brand_id
		LEFT JOIN product_category pc ON pc.id = p.category_id
		LEFT JOIN product_stock ps ON ps.product_id = p.id
		WHERE 1=1
	`
	
	args := []interface{}{}

	if search != "" {
		sql += " AND (p.product_name LIKE ? OR p.product_code LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if categoryID != "" {
		sql += " AND p.category_id = ?"
		args = append(args, categoryID)
	}

	if brandID != "" {
		sql += " AND p.brand_id = ?"
		args = append(args, brandID)
	}

	// Sorting
	orderClause := "p.product_name ASC"
	if sortBy == "product_name" {
		orderClause = "p.product_name"
	} else if sortBy == "product_code" {
		orderClause = "p.product_code"
	} else if sortBy == "brand_name" {
		orderClause = "b.brand_name"
	} else if sortBy == "product_category_name" {
		orderClause = "pc.product_category_name"
	} else if sortBy == "end_stock" {
		orderClause = "ps.end_stock"
	} else if sortBy == "last_buy_price" {
		orderClause = "ps.last_buy_price"
	} else if sortBy == "last_selling_price" {
		orderClause = "ps.last_selling_price"
	}

	if sortDir == "desc" {
		orderClause += " DESC"
	} else {
		orderClause += " ASC"
	}

	sql += " ORDER BY " + orderClause
	sql += " LIMIT ? OFFSET ?"
	args = append(args, limit, (page-1)*limit)

	err := r.db.Raw(sql, args...).Scan(&items).Error
	return items, total, err
}
