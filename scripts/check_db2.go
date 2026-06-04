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

	var res struct {
		Create_Table string `gorm:"column:Create Table"`
	}
	db.Raw("SHOW CREATE TABLE user_group_policies").Scan(&res)
	fmt.Println(res.Create_Table)
}
