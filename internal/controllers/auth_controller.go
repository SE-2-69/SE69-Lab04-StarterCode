package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"se69-lab04-starter/internal/config"
	"se69-lab04-starter/internal/dto"
	"se69-lab04-starter/internal/models"
	"se69-lab04-starter/internal/utils"
)

type AuthController struct {
	DB     *gorm.DB
	Config config.Config
}

// POST /sessions  (HTTP Basic: username and password)
func (ctl *AuthController) Login(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		ctl.unauthorized(c)
		return
	}

	var user models.User
	if ctl.DB.First(&user, "username = ?", username).Error != nil || !utils.CheckPassword(password, user.PasswordHash) {
		ctl.unauthorized(c)
		return
	}

	token, err := utils.CreateToken(user.Username, user.Role, ctl.Config.JWTSecret, ctl.Config.TokenTTL)
	if err != nil {
		utils.WriteError(c, http.StatusInternalServerError, "An unexpected error occurred")
		return
	}
	c.JSON(http.StatusCreated, dto.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: int(ctl.Config.TokenTTL.Seconds()),
	})
}

func (ctl *AuthController) unauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", `Basic realm="equipment-api"`)
	utils.WriteError(c, http.StatusUnauthorized, "Invalid username or password")
}
