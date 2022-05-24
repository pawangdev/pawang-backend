package models

import "time"

type User struct {
	ID           uint          `json:"id" form:"id"`
	Name         string        `json:"name" form:"name"`
	Email        string        `json:"email" form:"email"`
	Password     string        `json:"-" form:"password"`
	Phone        string        `json:"phone" form:"phone"`
	Gender       string        `json:"gender" form:"gender"`
	CreatedAt    time.Time     `json:"created_at" form:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" form:"updated_at"`
	Categories   []Category    `json:"-"`
	Transactions []Transaction `json:"-"`
	Wallets      []Wallet      `json:"-"`
}
