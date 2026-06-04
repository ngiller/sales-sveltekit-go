package main

import (
	"backend/config"
	"backend/internal/models"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func pointerToString(s string) *string {
	return &s
}

func pointerToUint(u uint) *uint {
	return &u
}

func main() {
	godotenv.Load("../backend/.env")
	db := config.InitDB()

	quotationID := "00012025100040"

	fmt.Println("=== STARTING END-TO-END FOLLOWUP DATABASE INTEGRATION TEST ===")
	fmt.Printf("Target Quotation ID: %s\n\n", quotationID)

	// 1. Calculate the next line_id sequence number
	fmt.Println("[Step 1] Calculating next sequence number (line_id)...")
	var maxLine int
	err := db.Model(&models.QuotationFollowup{}).
		Where("id = ?", quotationID).
		Select("COALESCE(MAX(line_id), 0)").
		Row().Scan(&maxLine)
	if err != nil {
		log.Fatalf("FAIL: Failed to calculate sequence number: %v", err)
	}
	nextLineID := maxLine + 1
	fmt.Printf("SUCCESS: Next sequence number (line_id) is: %d\n\n", nextLineID)

	// 2. Create the follow-up record
	fmt.Println("[Step 2] Inserting new follow-up record...")
	now := time.Now()
	followup := models.QuotationFollowup{
		ID:           quotationID,
		LineID:       nextLineID,
		PropertyID:   1,
		Notes:        pointerToString("Test integration followup after collation and PK fix"),
		Status:       pointerToUint(1),
		Progress:     pointerToUint(2),
		FollowupDate: &now,
		NextFollowup: &now,
	}
	
	err = db.Create(&followup).Error
	if err != nil {
		log.Fatalf("FAIL: Failed to insert follow-up record: %v", err)
	}
	fmt.Println("SUCCESS: Follow-up record inserted perfectly with no errors!")
	fmt.Printf("Record Data: ID (Quotation) = %s, LineID = %d, Notes = %s\n\n", followup.ID, followup.LineID, *followup.Notes)

	// 3. Query the inserted record to verify retrieval
	fmt.Println("[Step 3] Fetching inserted record back from database...")
	var retrieved models.QuotationFollowup
	err = db.Preload("StatusInfo").Preload("ProgressInfo").
		Where("id = ? AND line_id = ?", quotationID, nextLineID).
		First(&retrieved).Error
	if err != nil {
		log.Fatalf("FAIL: Failed to fetch follow-up record: %v", err)
	}
	fmt.Println("SUCCESS: Follow-up record retrieved perfectly with no errors!")
	fmt.Printf("Retrieved Data: ID = %s, LineID = %d, Notes = %s\n\n", retrieved.ID, retrieved.LineID, *retrieved.Notes)

	// 4. Update the record
	fmt.Println("[Step 4] Updating follow-up record...")
	retrieved.Notes = pointerToString("Test integration followup updated successfully")
	err = db.Save(&retrieved).Error
	if err != nil {
		log.Fatalf("FAIL: Failed to update follow-up record: %v", err)
	}
	fmt.Println("SUCCESS: Follow-up record updated perfectly with no errors!")

	// 5. Delete the record
	fmt.Println("[Step 5] Deleting follow-up record...")
	err = db.Where("id = ? AND line_id = ?", quotationID, nextLineID).Delete(&models.QuotationFollowup{}).Error
	if err != nil {
		log.Fatalf("FAIL: Failed to delete follow-up record: %v", err)
	}
	fmt.Println("SUCCESS: Follow-up record deleted perfectly with no errors!")
	fmt.Println("\n=== ALL FOLLOWUP INTEGRATION TESTS PASSED EXCELLENTLY ===")
}
