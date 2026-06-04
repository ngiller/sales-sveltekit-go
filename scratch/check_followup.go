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

	// Check if table exists
	hasTable := db.Migrator().HasTable("quotation_followup")
	fmt.Printf("Has quotation_followup table: %v\n", hasTable)

	if hasTable {
		// Print columns
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(struct{}{})
		columnTypes, err := db.Migrator().ColumnTypes("quotation_followup")
		if err != nil {
			log.Fatal(err)
		}
		for _, col := range columnTypes {
			dbType, _ := col.DatabaseTypeName()
			nullable, _ := col.Nullable()
			fmt.Printf("Col: %s, Type: %s, Nullable: %v\n", col.Name(), dbType, nullable)
		}
	} else {
		fmt.Println("Table does not exist. We need to create it.")
	}
}
