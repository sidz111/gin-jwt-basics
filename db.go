package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	db_username := "root"
	db_pass := "root"
	db_host := "localhost"
	db_port := "3303"
	db_name := "jwt_db"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_username, db_pass, db_host, db_port, db_name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	fmt.Println("Database connection established successfully.")
	return nil
}
