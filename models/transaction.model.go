package models

import (
	"time"
)

type Transaction struct {
	ID          uint      `json:"id" form:"id"`
	Amount      uint64    `json:"amount" form:"amount" gorm:"type:bigint;not null"`
	CategoryID  uint      `json:"category_id" form:"category_id" gorm:"not null"`
	WalletID    uint      `json:"wallet_id" form:"wallet_id" gorm:"not null"`
	Description string    `json:"description" form:"description" gorm:"type:text"`
	ImageUrl    string    `json:"image_url" form:"image_url" gorm:"type:varchar(255)"`
	Type        string    `json:"type" form:"type" gorm:"type:enum('income','outcome')"`
	Date        time.Time `json:"date" form:"date" gorm:"not null"`
	UserID      uint      `json:"user_id" form:"user_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" form:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" form:"updated_at"`
	Category    Category  `json:"category" gorm:"constraint:OnDelete:CASCADE;"`
	Wallet      Wallet    `json:"wallet" gorm:"constraint:OnDelete:CASCADE;"`
	User        User      `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}
