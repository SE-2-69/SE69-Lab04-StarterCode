package dto

type LoanRequest struct {
	EquipmentID uint `json:"equipmentId" binding:"required"`
}
