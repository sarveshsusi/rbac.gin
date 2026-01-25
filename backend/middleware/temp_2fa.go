package middleware

import (
	

	"github.com/gin-gonic/gin"
	

	"rbac/config"
	"rbac/utils"
)

func Temp2FAMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		raw := c.GetHeader("X-2FA-Token")
		if raw == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "missing 2fa token",
			})
			return
		}

		claims, err := utils.Parse2FAToken(
			raw,
			cfg.JWT.AccessSecret,
		)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid or expired 2fa session",
			})
			return
		}

		// 🔐 Attach user identity safely
		c.Set("two_fa_user_id", claims.UserID)

		c.Next()
	}
}
