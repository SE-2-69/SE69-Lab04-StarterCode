package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportController struct {
	DB *gorm.DB
}

type usageRow struct {
	Category string `json:"category"`
	Total    int    `json:"total"`
}

// GET /reports/usage  (teacher)
func (ctl *ReportController) Usage(c *gin.Context) {
	query := `SELECT category, COUNT(*) AS total FROM usage_reports GROUP BY category`
	rows := []usageRow{}
	if err := ctl.DB.Raw(query).Scan(&rows).Error; err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("database error: %v\nquery: %s\n", err, query))
		return
	}
	c.JSON(http.StatusOK, rows)
}
