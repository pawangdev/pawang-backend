package config

import (
	"fmt"
	"pawang-backend/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Database() (*gorm.DB, error) {
	DB_HOST := GetEnv("DB_HOST")
	DB_PORT := GetEnv("DB_PORT")
	DB_DATABASE := GetEnv("DB_DATABASE")
	DB_USERNAME := GetEnv("DB_USERNAME")
	DB_PASSWORD := GetEnv("DB_PASSWORD")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		DB_USERNAME,
		DB_PASSWORD,
		DB_HOST,
		DB_PORT,
		DB_DATABASE,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}

	db.AutoMigrate(entity.User{}, entity.Wallet{}, entity.Category{}, entity.Transaction{}, entity.SubCategory{}, entity.UserResetPassword{}, entity.UserNotification{})

	return db, nil
}
