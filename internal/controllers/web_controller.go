package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"se69-lab04-starter/internal/middleware"
	"se69-lab04-starter/internal/models"
	"se69-lab04-starter/internal/utils"
)

// routes for the web app in the browser: login with a form, then use a session cookie
type WebController struct {
	DB *gorm.DB
}

const sessionTTL = 30 * time.Minute

// POST /web/sessions  (form: username, password)
func (ctl *WebController) Login(c *gin.Context) {
	var user models.User
	if ctl.DB.First(&user, "username = ?", c.PostForm("username")).Error != nil ||
		!utils.CheckPassword(c.PostForm("password"), user.PasswordHash) {
		utils.WriteError(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	session := models.Session{ID: utils.RandomID(), Username: user.Username, ExpiresAt: time.Now().Add(sessionTTL)}
	if err := ctl.DB.Create(&session).Error; err != nil {
		utils.WriteError(c, http.StatusInternalServerError, "An unexpected error occurred")
		return
	}

	// SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
	c.SetCookie(middleware.SessionCookieName, session.ID, int(sessionTTL.Seconds()), "/web", "", false, false)
	c.JSON(http.StatusCreated, gin.H{"message": "Logged in as " + user.Username})
}

// DELETE /web/sessions  (log out)
func (ctl *WebController) Logout(c *gin.Context) {
	c.SetCookie(middleware.SessionCookieName, "", -1, "/web", "", false, false)
	c.Status(http.StatusNoContent)
}

// GET /web/loans  (my loans)
func (ctl *WebController) MyLoans(c *gin.Context) {
	list := []models.Loan{}
	if err := ctl.DB.Where("borrower = ?", middleware.GetClaims(c).Subject).Order("id").Find(&list).Error; err != nil {
		utils.WriteError(c, http.StatusInternalServerError, "An unexpected error occurred")
		return
	}
	c.JSON(http.StatusOK, list)
}

// POST /web/loans  (form: equipmentId)
func (ctl *WebController) Borrow(c *gin.Context) {
	equipmentID, err := strconv.ParseUint(c.PostForm("equipmentId"), 10, 64)
	if err != nil {
		utils.WriteError(c, http.StatusBadRequest, "equipmentId is required")
		return
	}
	borrow(c, ctl.DB, uint(equipmentID))
}
