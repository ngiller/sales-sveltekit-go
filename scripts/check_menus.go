package main

import (
	"backend/config"
	"backend/internal/models"
	"fmt"
	"log"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	var items []models.MasterTableAccess
	db.Find(&items)

	fmt.Println("ID | Name | MenuName")
	fmt.Println("--------------------")
	for _, item := range items {
		fmt.Printf("%d | %s | %s\n", item.ID, item.Name, item.MenuName)
	}
}
