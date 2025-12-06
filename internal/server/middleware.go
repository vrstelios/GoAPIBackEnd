package server

import (
	"GoAPIBackEnd/config"
	"GoAPIBackEnd/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

var JwtSecret []byte

func AuthJWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tokenString string

		authHeader := ctx.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := ctx.Cookie("Authorization")
			if err != nil || cookie == "" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
				return
			}
			tokenString = cookie
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return JwtSecret, nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		ctx.Set("username", claims["username"])
		ctx.Set("role", claims["role"])
		ctx.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString("username")
		user, ok := models.Users[username]
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		for _, role := range allowedRoles {
			if user.Role == role {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	}
}

func init() {
	config.LoadEnvVariables()
	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
}
