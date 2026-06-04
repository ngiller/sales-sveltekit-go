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

	id := "202605094971328"
	db.Table("quotation_master").Where("id = ?", id).Update("default_quot", false)
	db.Table("quotation_master").Where("id = ? AND rev_id = ?", id, 1).Update("default_quot", true)
	fmt.Println("Fixed default status for 202605094971328 Rev 1")
}
