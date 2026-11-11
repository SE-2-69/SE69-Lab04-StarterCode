package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"se69-lab04-starter/internal/models"
	"se69-lab04-starter/internal/utils"
)

const SessionCookieName = "SESSIONID"

// check the session cookie (used by the /web routes instead of a Bearer token)
func RequireSession(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(SessionCookieName)
		if err != nil {
			utils.WriteError(c, http.StatusUnauthorized, "Please log in first")
			return
		}

		var session models.Session
		var user models.User
		if db.First(&session, "id = ?", sessionID).Error != nil ||
			time.Now().After(session.ExpiresAt) ||
			db.First(&user, "username = ?", session.Username).Error != nil {
			utils.WriteError(c, http.StatusUnauthorized, "Session is invalid or expired")
			return
		}

		claims := &utils.Claims{Role: user.Role}
		claims.Subject = user.Username
		c.Set(claimsKey, claims)
		c.Next()
	}
}
