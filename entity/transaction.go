package entity

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID            int
	Amount        int           `gorm:"type:bigint(20);not null"`
	CategoryID    int           `gorm:"not null"`
	SubCategoryID sql.NullInt32 `gorm:"default:null"`
	WalletID      int           `gorm:"not null"`
	Type          string        `gorm:"type:enum('income', 'outcome');not null"`
	Description   string        `gorm:"type:text"`
	ImageUrl      string        `gorm:"type:varchar(255);not null"`
	Date          time.Time     `gorm:"not null"`
	UserID        int           `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Category      Category
	SubCategory   SubCategory
	Wallet        Wallet
	User          User
}
