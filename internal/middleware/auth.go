package middleware

import (
	"GoAPIBackEnd/internal/helpers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				token = parts[1]
			}
		}

		if token == "" {
			cookieToken, err := c.Cookie("Authorization")
			if err == nil && cookieToken != "" {
				token = cookieToken
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token required", "code": "UNAUTHORIZED"})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(token)
		if err != nil {
			log.Printf("Token validation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token", "code": "INVALID_TOKEN"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsInterface, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required", "code": "UNAUTHORIZED"})
			c.Abort()
			return
		}

		claims, ok := claimsInterface.(*helpers.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication data", "code": "INVALID_CLAIMS"})
			c.Abort()
			return
		}

		userHasAccess := false
		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(claims.Role, allowedRole) {
				userHasAccess = true
				break
			}
		}

		if !userHasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "FORBIDDEN", "user_role": claims.Role, "required_roles": allowedRoles})
			c.Abort()
			return
		}

		c.Next()
	}
}
