package middleware

import (
	"GoAPIBackEnd/internal/apperrors"
	"GoAPIBackEnd/internal/helpers"
	"github.com/gin-gonic/gin"
	"strings"
)

func Authenticate(provider TokenProvider) gin.HandlerFunc {
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
			apperrors.GetAPIError(c, nil, 0, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := provider.Validate(token)
		if err != nil {
			apperrors.GetAPIError(c, nil, 0, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		c.Set(helpers.CtxUserId, claims.UserId)
		c.Set(helpers.CtxEmail, claims.Email)
		c.Set(helpers.CtxRole, claims.Role)

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := helpers.GetUserRole(c)
		if role == "" {
			apperrors.GetAPIError(c, nil, 0, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		userHasAccess := false
		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(role, allowedRole) {
				userHasAccess = true
				break
			}
		}

		if !userHasAccess {
			apperrors.GetAPIError(c, nil, 0, apperrors.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
