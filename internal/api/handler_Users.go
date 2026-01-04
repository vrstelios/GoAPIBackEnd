package api

import (
	"GoAPIBackEnd/config"
	"GoAPIBackEnd/internal/apperrors"
	"GoAPIBackEnd/internal/auth"
	"GoAPIBackEnd/internal/database"
	"GoAPIBackEnd/internal/type/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// @Summary Signup a new user.
// @Description Creates a new user account with username and password
// @Tags auth
// @Produce json
// @Param user body models.Users true "User registration data"
// @Success 200 {object} models.APIError
// @Failure 400 {object} models.APIError
// @Failure 406 {object} models.APIError
// @Failure 409 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /auth/signup [post]
func Signup(ctx *gin.Context) {
	var dbConn = database.Conn

	userBody := models.Users{}
	err := ctx.ShouldBindJSON(&userBody)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	if len(userBody.Role) == 0 {
		userBody.Role = "user"
	}

	if len(userBody.Name) == 0 && len(userBody.Email) == 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotAcceptable, "INVALID_CREDENTIALS", fmt.Errorf("Invalid username/password")))
		return
	}
	if len(userBody.Password) < 8 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotAcceptable, "INVALID_PASSWORD", fmt.Errorf("Password too short")))
		return
	}

	user, err := database.GetUser(dbConn, "", userBody.Name)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}
	if len(user) > 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusConflict, "USER_EXISTS", fmt.Errorf("User already exists")))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userBody.Password), 10)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "HASH_ERROR", err))
		return
	}

	nUser := models.Users{
		Id:       uuid.New().String(),
		Name:     userBody.Name,
		Password: string(hashedPassword),
		Email:    userBody.Email,
		Role:     userBody.Role,
	}

	err = database.PostUser(dbConn, nUser)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "INVALID_REQUEST", err))
		return
	}
	apperrors.GetAPIError(ctx, gin.H{"message": "User registered successfully!"}, http.StatusOK, nil)
}

// @Summary Login user
// @Description Authenticates user and returns session cookies
// @Tags auth
// @Produce json
// @Param user body models.Users true "User registration data"
// @Success 200 {object} models.APIError
// @Failure 400 {object} models.APIError
// @Failure 401 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /auth/login [post]
func Login(ctx *gin.Context) {
	var dbConn = database.Conn

	userPayload := models.Users{}
	if err := ctx.ShouldBindJSON(&userPayload); err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	user, err := database.GetUser(dbConn, "", userPayload.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusUnauthorized, "INVALID_CREDENTIALS", fmt.Errorf("Invalid username or password")))
		return
	}
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user[0].Password), []byte(userPayload.Password))
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusUnauthorized, "INVALID_CREDENTIALS", fmt.Errorf("Invalid username or password")))
		return
	}

	// Generate a jwt token
	expHours := config.GetConfig().JWT.ExpirationHours
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userPayload.Name,
		"role":     user[0].Role,
		"exp":      time.Now().Add(time.Duration(expHours) * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(auth.JwtSecret)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "TOKEN_GENERATION_ERROR", fmt.Errorf("Failed to generate token: %v", err)))
		return
	}

	// Send it back
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, expHours*3600, "/", "", false, true)

	apperrors.GetAPIError(ctx, gin.H{}, http.StatusOK, nil)
}

// @Summary	Logout user.
// @Tags		auth
// @Produce	json
// @Success	200		{array}		models.APIError
// @Failure	500		{object}	models.APIError
// @Router		/auth/logout [get]
func Logout(ctx *gin.Context) {
	// Delete the JWT cookie
	ctx.SetCookie("Authorization", "", -1, "/", "", false, true)

	apperrors.GetAPIError(ctx, gin.H{"message": "Logout successfully!"}, http.StatusOK, nil)
}
