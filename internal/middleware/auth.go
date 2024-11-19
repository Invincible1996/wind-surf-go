package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"wind-surf-go/internal/utils"
)

// AuthMiddleware verifies the JWT token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "missing authorization header")
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "invalid authorization header format")
			c.Abort()
			return
		}

		// Verify token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.UnauthorizedResponse(c, "invalid token")
			c.Abort()
			return
		}

		// Store user info in context for later use
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
