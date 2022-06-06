package models

import "time"

type Wallet struct {
	ID        uint      `json:"id" form:"id"`
	Name      string    `json:"name" form:"name" gorm:"type:varchar(100);not null"`
	UserID    uint      `json:"user_id" form:"user_id" gorm:"not null"`
	Balance   uint64    `json:"balance" form:"balance" gorm:"type:bigint;not null"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
	User      User      `json:"-"`
}
