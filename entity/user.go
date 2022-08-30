package entity

import "time"

type User struct {
	ID            int
	Name          string `gorm:"type:varchar(50);not null;"`
	Email         string `gorm:"type:varchar(50);not null;unique"`
	Password      string `gorm:"type:varchar(255);not null;"`
	Phone         string `gorm:"type:varchar(15);not null;"`
	Gender        string `gorm:"type:enum('male', 'female');not null;"`
	ImageProfile  string `gorm:"type:varchar(255)"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Wallets       []Wallet
	SubCategories []SubCategory
	Transactions  []Transaction
}
