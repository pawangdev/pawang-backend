package models

import "time"

type User struct {
	ID           uint          `json:"id" form:"id"`
	Name         string        `json:"name" form:"name" gorm:"type:varchar(100);not null"`
	Email        string        `json:"email" form:"email" gorm:"type:varchar(100);unique;not null"`
	Password     string        `json:"-" form:"password" gorm:"type:varchar(255);not null"`
	Phone        string        `json:"phone" form:"phone" gorm:"type:varchar(15);not null"`
	Gender       string        `json:"gender" form:"gender" gorm:"type:varchar(10);not null"`
	CreatedAt    time.Time     `json:"created_at" form:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" form:"updated_at"`
	Categories   []Category    `json:"-"`
	Transactions []Transaction `json:"-"`
	Wallets      []Wallet      `json:"-"`
}
