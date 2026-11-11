package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"se69-lab04-starter/internal/dto"
	"se69-lab04-starter/internal/models"
	"se69-lab04-starter/internal/utils"
)

type EquipmentController struct {
	DB *gorm.DB
}

// GET /equipment?category=board
func (ctl *EquipmentController) List(c *gin.Context) {
	list := []models.Equipment{}
	query := ctl.DB.Order("id")
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if err := query.Find(&list).Error; err != nil {
		utils.WriteError(c, http.StatusInternalServerError, "An unexpected error occurred")
		return
	}
	c.JSON(http.StatusOK, list)
}

// GET /equipment/:id
func (ctl *EquipmentController) Get(c *gin.Context) {
	equipment, ok := ctl.find(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, equipment)
}

// POST /equipment  (staff)
func (ctl *EquipmentController) Create(c *gin.Context) {
	var req dto.EquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, http.StatusBadRequest, "name and category are required")
		return
	}

	equipment := models.Equipment{Name: req.Name, Category: req.Category, Status: "available"}
	if err := ctl.DB.Create(&equipment).Error; err != nil {
		utils.WriteError(c, http.StatusInternalServerError, "An unexpected error occurred")
		return
	}
	c.Header("Location", fmt.Sprintf("/equipment/%d", equipment.ID))
	c.JSON(http.StatusCreated, equipment)
}

// PUT /equipment/:id  (staff)
func (ctl *EquipmentController) Update(c *gin.Context) {
	equipment, ok := ctl.find(c)
	if !ok {
		return
	}
	var req dto.EquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, http.StatusBadRequest, "name and category are required")
		return
	}

	equipment.Name = req.Name
	equipment.Category = req.Category
	ctl.DB.Save(&equipment)
	c.JSON(http.StatusOK, equipment)
}

// GET /deleteEquipment/:id  (staff)
func (ctl *EquipmentController) Delete(c *gin.Context) {
	id, _ := utils.ParamID(c)
	var equipment models.Equipment
	if ctl.DB.First(&equipment, id).Error != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Not found"})
		return
	}
	ctl.DB.Delete(&equipment)
	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

// find equipment by :id, or send 404
func (ctl *EquipmentController) find(c *gin.Context) (models.Equipment, bool) {
	var equipment models.Equipment
	id, ok := utils.ParamID(c)
	if !ok || ctl.DB.First(&equipment, id).Error != nil {
		utils.WriteError(c, http.StatusNotFound, fmt.Sprintf("Equipment %s doesn't exist", c.Param("id")))
		return equipment, false
	}
	return equipment, true
}
