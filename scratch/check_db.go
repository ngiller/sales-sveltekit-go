package main

import (
	"fmt"
	"time"
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Quotation struct {
	ID           string     `gorm:"column:id"`
	QuotationID  string     `gorm:"column:quotation_id"`
	GrandTotal   float64    `gorm:"column:grand_total"`
	UserCreated  uint       `gorm:"column:user_created"`
	NextFollowup *time.Time `gorm:"column:next_followup"`
}

type QuotationMaster struct {
	ID          string  `gorm:"column:id"`
	RevID       int     `gorm:"column:rev_id"`
	GrandTotal  float64 `gorm:"column:grand_total"`
	DefaultQuot bool    `gorm:"column:default_quot"`
}

func main() {
	dsn := "root:Pass@w0rd@tcp(localhost:3306)/magnum_sales_svelte_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	id := "202605094971328"

	var q Quotation
	db.Table("quotation").Where("id = ?", id).First(&q)
	fmt.Printf("Quotation Table: ID=%s, QuotID=%s, GrandTotal=%.2f, CreatedBy=%d, NextFollowup=%v\n", q.ID, q.QuotationID, q.GrandTotal, q.UserCreated, q.NextFollowup)

	var masters []QuotationMaster
	db.Table("quotation_master").Where("id = ?", id).Find(&masters)
	fmt.Println("Quotation Master Revisions:")
	for _, m := range masters {
		fmt.Printf("  RevID=%d, GrandTotal=%.2f, Default=%v\n", m.RevID, m.GrandTotal, m.DefaultQuot)
	}
}
