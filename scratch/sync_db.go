package main

import (
	"fmt"
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type QuotationMaster struct {
	ID          string   `gorm:"column:id"`
	RevID       int      `gorm:"column:rev_id"`
	Total       *float64 `gorm:"column:total"`
	Disc        *float64 `gorm:"column:disc"`
	Tax         *float64 `gorm:"column:tax"`
	TaxValue    float64  `gorm:"column:tax_value"`
	PPh         float64  `gorm:"column:pph"`
	PPhValue    float64  `gorm:"column:pph_value"`
	GrandTotal  *float64 `gorm:"column:grand_total"`
	HppTotal    *float64 `gorm:"column:hpp_total"`
	Profit      float64  `gorm:"column:profit"`
	Margin      *float64 `gorm:"column:margin"`
	ProfitValue *float64 `gorm:"column:profit_value"`
	Notes       *string  `gorm:"column:notes"`
	DefaultQuot bool     `gorm:"column:default_quot"`
}

func main() {
	dsn := "root:Pass@w0rd@tcp(localhost:3306)/magnum_sales_svelte_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	id := "202605094971328"

	var master QuotationMaster
	if err := db.Table("quotation_master").Where("id = ? AND default_quot = ?", id, true).First(&master).Error; err != nil {
		log.Fatalf("No default revision found for %s: %v", id, err)
	}

	fmt.Printf("Syncing ID %s with Default Rev %d\n", id, master.RevID)

	updates := map[string]interface{}{
		"total":        master.Total,
		"disc":         master.Disc,
		"tax":          master.Tax,
		"tax_value":    master.TaxValue,
		"pph":          master.PPh,
		"pph_value":    master.PPhValue,
		"grand_total":  master.GrandTotal,
		"hpp_total":    master.HppTotal,
		"profit":       master.Profit,
		"margin":       master.Margin,
		"profit_value": master.ProfitValue,
		"notes":        master.Notes,
	}

	if err := db.Table("quotation").Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sync successful!")
}
