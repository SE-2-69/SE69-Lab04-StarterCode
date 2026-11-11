package models

type User struct {
	Username     string `gorm:"primaryKey"`
	Role         string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}
