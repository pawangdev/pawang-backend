package models

import "time"

type Category struct {
	ID                uint            `json:"id" form:"id"`
	Name              string          `json:"name" form:"name" gorm:"type:varchar(100);not null"`
	IconUrl           string          `json:"icon_url" form:"icon_url" gorm:"type:varchar(255);not null"`
	UserID            uint            `json:"user_id" form:"user_id"`
	TransactionTypeID uint            `json:"transaction_type_id" form:"transaction_type_id"`
	CreatedAt         time.Time       `json:"created_at" form:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" form:"updated_at"`
	User              User            `json:"-"`
	TransactionType   TransactionType `json:"transaction_type"`
	Transactions      []Transaction   `json:"-"`
}
