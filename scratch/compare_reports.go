package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SalesChartsBySalesPersonItem struct {
	CustomerName  string  `json:"customer_name"`
	GrandTotal    float64 `json:"grand_total"`
	HppTotal      float64 `json:"hpp_total"`
	ProfitValue   float64 `json:"profit_value"`
	MarginPercent float64 `json:"margin_percent"`
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file, trying local dir")
		godotenv.Load(".env")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	salesIDStr := "114" // Wanda

	ranges := [][]string{
		{"2024-01-01", "2024-12-31"},
		{"2024-01-01", "2024-06-30"},
		{"2024-07-01", "2024-12-31"},
	}

	for _, r := range ranges {
		fromDate := r[0]
		toDate := r[1]
		fmt.Printf("\n--- Testing range: %s to %s ---\n", fromDate, toDate)

		var items []SalesChartsBySalesPersonItem
		if err := db.Table("quotation q").
			Select(`c.name AS customer_name, COALESCE(SUM(q.grand_total),0) AS grand_total, COALESCE(SUM(q.hpp_total),0) AS hpp_total, COALESCE(SUM(q.grand_total - q.hpp_total),0) AS profit_value, CASE WHEN SUM(q.grand_total) > 0 THEN ROUND((SUM(q.grand_total - q.hpp_total) / SUM(q.grand_total)) * 100, 2) ELSE 0 END AS margin_percent`).
			Joins("JOIN customer c ON c.id = q.customer_id").
			Where("q.sales_id = ?", salesIDStr).
			Where("q.status = 3 AND q.progress = 9").
			Where("q.quotation_date BETWEEN ? AND ?", fromDate, toDate).
			Group("c.id, c.name").
			Order("grand_total DESC").
			Find(&items).Error; err != nil {
			log.Fatalf("failed items: %v", err)
		}

		var chartSum float64
		for _, item := range items {
			chartSum += item.GrandTotal
		}
		fmt.Printf("PO Grid Grand Total: %f (items count: %d)\n", chartSum, len(items))
	}
}
