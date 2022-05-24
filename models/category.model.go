package models

import "time"

type Category struct {
	ID           uint          `json:"id" form:"id"`
	Name         string        `json:"name" form:"name"`
	IconUrl      string        `json:"icon_url" form:"icon_url"`
	UserID       uint          `json:"user_id" form:"user_id"`
	CreatedAt    time.Time     `json:"created_at" form:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" form:"updated_at"`
	User         User          `json:"-"`
	Transactions []Transaction `json:"-"`
}
