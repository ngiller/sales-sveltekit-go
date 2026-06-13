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
	// Uses LEFT JOIN to a derived table for last buy price (unit price, not total)
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
			COALESCE(lp.last_buy_unit_price, 0) as last_buy_price,
			COALESCE(ps.last_selling_price, '0') as last_selling_price,
			ps.last_buy_date,
			ps.unit_id
		FROM product p
		LEFT JOIN brand b ON b.id = p.brand_id
		LEFT JOIN product_category pc ON pc.id = p.category_id
		LEFT JOIN product_stock ps ON ps.product_id = p.id
		LEFT JOIN (
			SELECT product_id, total_price as last_buy_unit_price
			FROM (
				SELECT *,
					ROW_NUMBER() OVER (PARTITION BY product_id ORDER BY tgl DESC) as rn
				FROM (
					SELECT ps2.product_id, pop.total_price, pop.quantity as qty, po.purchase_date as tgl
					FROM purchase_order_product pop
					JOIN purchase_order po ON po.id = pop.purchase_order_id AND po.status_receiving = 'CLOSE'
					JOIN product_stock ps2 ON ps2.id = pop.product_stock_id
					WHERE pop.quantity > 0

					UNION ALL

					SELECT ps2.product_id, dbp.total_price, dbp.quantity as qty, dbp.purchase_date as tgl
					FROM direct_buying_product dbp
					JOIN product_stock ps2 ON ps2.id = dbp.product_stock_id
					WHERE dbp.quantity > 0
				) combined_all
			) ranked_purchases
			WHERE rn = 1
		) lp ON lp.product_id = p.id
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
		orderClause = "last_buy_price"
	} else if sortBy == "last_selling_price" {
		orderClause = "ps.last_selling_price"
	} else if sortBy == "last_buy_date" {
		orderClause = "ps.last_buy_date"
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

func (r *StockRepository) GetTotalStockValue(search string, categoryID, brandID string) (float64, error) {
	var totalValue float64

	sql := `
		SELECT COALESCE(SUM(COALESCE(ps.end_stock, 0) * COALESCE(ps.last_selling_price, 0)), 0)
		FROM product p
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

	err := r.db.Raw(sql, args...).Scan(&totalValue).Error
	return totalValue, err
}

// FindPurchaseHistory returns the last 10 purchase transactions for a product.
// It queries both purchase_order and direct_buying tables from magnum_stock_db.
func (r *StockRepository) FindPurchaseHistory(productID uint) ([]map[string]interface{}, error) {
	sql := `
		SELECT sub.purchase_date, v.vendor_name, sub.qty, sub.price, sub.total
		FROM (
			SELECT 
				po.purchase_date,
				po.vendor_id,
				pop.quantity AS qty,
				pop.total_price AS price,
				ROUND(pop.total_price * pop.quantity, 0) AS total
			FROM purchase_order_product pop
			JOIN purchase_order po ON po.id = pop.purchase_order_id
			WHERE pop.product_stock_id IN (SELECT id FROM product_stock WHERE product_id = ?)
			  AND po.status_receiving = 'CLOSE'

			UNION ALL

			SELECT 
				dbp.purchase_date,
				db.vendor_id,
				dbp.quantity AS qty,
				dbp.total_price AS price,
				ROUND(dbp.total_price * dbp.quantity, 0) AS total
			FROM direct_buying_product dbp
			JOIN direct_buying db ON db.id = dbp.direct_buying_id
			WHERE dbp.product_stock_id IN (SELECT id FROM product_stock WHERE product_id = ?)
		) sub
		LEFT JOIN vendor v ON v.id = CAST(sub.vendor_id AS UNSIGNED)
		ORDER BY sub.purchase_date DESC
		LIMIT 10
	`

	var items []map[string]interface{}
	err := r.db.Raw(sql, productID, productID).Scan(&items).Error
	return items, err
}
