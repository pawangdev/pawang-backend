package entity

import "time"

type UserNotification struct {
	ID          int
	UserID      int    `gorm:"not null"`
	OnesignalID string `gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        User
}
