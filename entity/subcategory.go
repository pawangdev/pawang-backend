package entity

import "time"

type SubCategory struct {
	ID           int
	Name         string `gorm:"type:varchar(100);not null;unique"`
	CategoryID   int    `gorm:"not null"`
	UserID       int    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Category     Category
	Transactions []Transaction `gorm:"constraint:OnDelete:SET NULL;"`
	User         User
}
