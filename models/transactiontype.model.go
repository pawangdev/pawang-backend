package models

import "time"

type TransactionType struct {
	ID          uint          `json:"id" form:"id"`
	Name        string        `json:"name" form:"name" gorm:"type:varchar(100);not null"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Transaction []Transaction `json:"-"`
	Category    []Category    `json:"-"`
}
