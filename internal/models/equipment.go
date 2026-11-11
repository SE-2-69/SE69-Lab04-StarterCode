package models

type Equipment struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Category string `gorm:"not null" json:"category"`
	Status   string `gorm:"not null;default:available" json:"status"`
}
