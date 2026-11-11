package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"se69-lab04-starter/internal/dto"
	"se69-lab04-starter/internal/middleware"
	"se69-lab04-starter/internal/models"
	"se69-lab04-starter/internal/utils"
)

type LoanController struct {
	DB *gorm.DB
}

// POST /loans  (student)
func (ctl *LoanController) Create(c *gin.Context) {
	var req dto.LoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, http.StatusBadRequest, "equipmentId is required")
		return
	}
	borrow(c, ctl.DB, req.EquipmentID)
}

// GET /loans  (students see only their own loans)
func (ctl *LoanController) List(c *gin.Context) {
	claims := middleware.GetClaims(c)
	list := []models.Loan{}
	query := ctl.DB.Order("id")
	if claims.Role == "student" {
		query = query.Where("borrower = ?", claims.Subject)
	}
	if err := query.Find(&list).Error; err != nil {
		utils.WriteError(c, http.StatusInternalServerError, "An unexpected error occurred")
		return
	}
	c.JSON(http.StatusOK, list)
}

// GET /loans/:id
func (ctl *LoanController) Get(c *gin.Context) {
	var loan models.Loan
	id, ok := utils.ParamID(c)
	if !ok || ctl.DB.First(&loan, id).Error != nil {
		utils.WriteError(c, http.StatusNotFound, fmt.Sprintf("Loan %s doesn't exist", c.Param("id")))
		return
	}
	c.JSON(http.StatusOK, loan)
}

// borrow equipment for the logged-in user (used by POST /loans and POST /web/loans)
func borrow(c *gin.Context, db *gorm.DB, equipmentID uint) {
	var equipment models.Equipment
	if db.First(&equipment, equipmentID).Error != nil {
		utils.WriteError(c, http.StatusBadRequest, "equipmentId doesn't match any equipment")
		return
	}
	if equipment.Status != "available" {
		utils.WriteError(c, http.StatusConflict, "Equipment is already borrowed")
		return
	}

	loan := models.Loan{EquipmentID: equipment.ID, Borrower: middleware.GetClaims(c).Subject}
	if err := db.Create(&loan).Error; err != nil {
		utils.WriteError(c, http.StatusInternalServerError, "An unexpected error occurred")
		return
	}
	db.Model(&equipment).Update("status", "borrowed")

	c.Header("Location", fmt.Sprintf("/loans/%d", loan.ID))
	c.JSON(http.StatusCreated, loan)
}
