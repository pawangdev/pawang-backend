package models

import "time"

type User struct {
	ID           uint          `json:"id" form:"id"`
	Name         string        `json:"name" form:"name" gorm:"type:varchar(100);not null"`
	Email        string        `json:"email" form:"email" gorm:"type:varchar(100);unique;not null"`
	Password     string        `json:"-" form:"password" gorm:"type:varchar(255);not null"`
	Phone        string        `json:"phone" form:"phone" gorm:"type:varchar(15);"`
	Gender       string        `json:"gender" form:"gender" gorm:"type:varchar(10);"`
	ImageProfile string        `json:"image_profile" form:"image_profile" gorm:"type:varchar(255)"`
	CreatedAt    time.Time     `json:"created_at" form:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" form:"updated_at"`
	Categories   []Category    `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	Transactions []Transaction `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	Wallets      []Wallet      `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}
