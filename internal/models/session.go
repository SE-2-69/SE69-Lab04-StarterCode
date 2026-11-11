package models

import "time"

// server-side session for the /web routes (the cookie only holds the ID)
type Session struct {
	ID        string `gorm:"primaryKey"`
	Username  string `gorm:"not null"`
	ExpiresAt time.Time
}
