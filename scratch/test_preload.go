package main

import (
	"backend/config"
	"backend/internal/models"
	"fmt"
	"log"
)

func main() {
	// Initialize config
	config.InitConfig()
	db := config.InitDB().Debug()

	var details []models.QuotationDetail
	err := db.Preload("Unit").Limit(5).Find(&details).Error
	if err != nil {
		log.Fatalf("Failed to query quotation details with preload: %v", err)
	}

	fmt.Printf("Successfully loaded %d quotation details:\n", len(details))
	for _, d := range details {
		unitName := "nil"
		if d.Unit != nil {
			unitName = d.Unit.Name
		}
		var unitIDVal interface{} = "nil"
		if d.UnitID != nil {
			unitIDVal = *d.UnitID
		}
		fmt.Printf("  Detail ID=%s, Line=%d, UnitID=%v, UnitName=%s\n", d.ID, d.Line, unitIDVal, unitName)
	}
}
