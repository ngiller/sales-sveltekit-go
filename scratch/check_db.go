package main

import (
	"fmt"
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:Pass@w0rd@tcp(localhost:3306)/magnum_sales_svelte_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	type Result struct {
		Status      int
		Progress    int
		Count       int64
		GrandTotal  float64
	}

	var results []Result
	db.Table("quotation").
		Select("status, progress, COUNT(*) as count, SUM(grand_total) as grand_total").
		Group("status, progress").
		Find(&results)

	fmt.Println("--- Grouped by Status & Progress ---")
	for _, r := range results {
		fmt.Printf("Status: %d, Progress: %d, Count: %d, GrandTotal: %.2f\n", r.Status, r.Progress, r.Count, r.GrandTotal)
	}
}
