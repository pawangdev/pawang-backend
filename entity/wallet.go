package entity

import "time"

type Wallet struct {
	ID           int
	Name         string `gorm:"type:varchar(100);not null"`
	UserID       int    `gorm:"not null"`
	Balance      int    `gorm:"type:bigint(20)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	User         User
	Transactions []Transaction `gorm:"constraint:OnDelete:CASCADE;"`
}
