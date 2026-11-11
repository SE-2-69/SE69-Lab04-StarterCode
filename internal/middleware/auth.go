package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"se69-lab04-starter/internal/utils"
)

const claimsKey = "claims"

// check the Bearer token, then check the user's role (no roles = any logged-in user)
func RequireAuth(secret string, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.Header("WWW-Authenticate", `Bearer realm="equipment-api"`)
			utils.WriteError(c, http.StatusUnauthorized, "Missing, invalid, or expired token")
			return
		}

		claims, err := utils.VerifyToken(strings.TrimPrefix(header, "Bearer "), secret)
		if err != nil {
			c.Header("WWW-Authenticate", `Bearer realm="equipment-api", error="invalid_token"`)
			utils.WriteError(c, http.StatusUnauthorized, "Missing, invalid, or expired token")
			return
		}

		if len(roles) > 0 && !slices.Contains(roles, claims.Role) {
			utils.WriteError(c, http.StatusForbidden, "Your role is not allowed to perform this operation")
			return
		}

		c.Set(claimsKey, claims)
		c.Next()
	}
}

// get the logged-in user's claims inside a controller
func GetClaims(c *gin.Context) *utils.Claims {
	return c.MustGet(claimsKey).(*utils.Claims)
}
