package models

import (
	"time"
)

type Transaction struct {
	ID          uint      `json:"id" form:"id"`
	Amount      uint64    `json:"amount" form:"amount"`
	CategoryID  uint      `json:"category_id" form:"category_id"`
	WalletID    uint      `json:"wallet_id" form:"wallet_id"`
	Type        string    `json:"type" form:"type"`
	Description string    `json:"description" form:"description"`
	ImageUrl    string    `json:"image_url" form:"image_url"`
	Date        time.Time `json:"date" form:"date"`
	UserID      uint      `json:"user_id" form:"user_id"`
	CreatedAt   time.Time `json:"created_at" form:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" form:"updated_at"`
	Category    Category  `json:"-"`
	Wallet      Wallet    `json:"-"`
	User        User      `json:"-"`
}
