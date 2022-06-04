package models

import "time"

type Wallet struct {
	ID        uint      `json:"id" form:"id"`
	Name      string    `json:"name" form:"name"`
	UserID    uint      `json:"user_id" form:"user_id"`
	Balance   uint64    `json:"balance" form:"balance"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
	User      User      `json:"-"`
}
