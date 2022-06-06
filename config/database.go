package config

import (
	"fmt"
	"log"
	"pawang-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	username := GetEnv("DB_USERNAME")
	password := GetEnv("DB_PASSWORD")
	host := GetEnv("DB_HOST")
	database := GetEnv("DB_DATABASE")

	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	db.AutoMigrate(models.Transaction{}, models.User{}, models.Category{}, models.TransactionType{}, models.Wallet{})

	return db
}
