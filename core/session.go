package core

import (
	"GoAPIBackEnd/model"
	"errors"
	"github.com/gin-gonic/gin"
	"net/url"
)

var AuthErr = errors.New("Unauthorized")

func Authorize(ctx *gin.Context) error {
	// Get session token from cookie
	st, err := ctx.Cookie("session_token")
	if err != nil || st == "" {
		return AuthErr
	}
	stDecoded, _ := url.QueryUnescape(st)

	// Get CSRF token from header
	csrf := ctx.GetHeader("X-CSRF-Token")
	if csrf == "" {
		return AuthErr
	}
	csrfDecoded, _ := url.QueryUnescape(csrf)

	// Find user by session token
	for username, user := range model.Users {
		if user.SessionToken == stDecoded && user.CSRFToken == csrfDecoded {
			ctx.Set("username", username)
			return nil
		}
	}

	return AuthErr
}
