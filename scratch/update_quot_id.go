package main

import (
	"fmt"
	"log"
	"strings"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:Pass@w0rd@tcp(localhost:3306)/magnum_sales_svelte_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	id := "202605094971328"
	var currentID string
	var userID uint
	db.Table("quotation").Select("quotation_id, user_created").Where("id = ?", id).Row().Scan(&currentID, &userID)

	var inisial string
	db.Table("users").Select("inisial").Where("id = ?", userID).Scan(&inisial)

	if inisial == "" {
		inisial = "USR"
	}

	newID := strings.Replace(currentID, "/USR", "/"+inisial, 1)
	fmt.Printf("Updating Quotation ID: %s -> %s\n", currentID, newID)

	db.Table("quotation").Where("id = ?", id).Update("quotation_id", newID)
	fmt.Println("Update successful!")
}
