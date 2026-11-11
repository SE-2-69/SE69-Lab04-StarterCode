package dto

// binding:"required" = Gin returns an error if the field is missing
type EquipmentRequest struct {
	Name     string `json:"name" binding:"required"`
	Category string `json:"category" binding:"required"`
}
