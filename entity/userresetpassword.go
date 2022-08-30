package entity

import "time"

type UserResetPassword struct {
	ID        int
	Email     string `gorm:"type:varchar(50);not null"`
	Token     string `gorm:"type:varchar(8);not null"`
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
