package entity

import "time"

type Category struct {
	ID        int
	Name      string `gorm:"type:varchar(100);not null"`
	Icon      string `gorm:"type:text;not null"`
	Type      string `gorm:"type:enum('income', 'outcome');not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// SubCategories []SubCategory `gorm:"constraint:OnDelete:CASCADE;"`
}
