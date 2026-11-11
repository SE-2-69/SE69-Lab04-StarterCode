package models

import "time"

type Loan struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	EquipmentID uint      `gorm:"not null" json:"equipmentId"`
	Borrower    string    `gorm:"not null" json:"borrower"`
	BorrowedAt  time.Time `gorm:"autoCreateTime" json:"borrowedAt"`
}
