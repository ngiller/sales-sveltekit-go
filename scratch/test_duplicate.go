package main

import (
	"backend/internal/repository"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:Pass@w0rd@tcp(localhost:3306)/magnum_sales_svelte_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewQuotationRepository(db)

	oldID := "202605094971328"
	newID := repository.GenerateID()
	newQuotationID := "TEST/DUP/001"
	now := time.Now()
	userID := uint(1)
	userInisial := "SYS"

	err = repo.CreateRevision(oldID, newID, newQuotationID, now, userID, userInisial, "[DUPLICATE] ")
	if err != nil {
		fmt.Printf("Error duplicating: %v\n", err)
	} else {
		fmt.Println("Duplication successful!")
	}
}
