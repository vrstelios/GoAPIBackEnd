package api

import (
	"GoAPIBackEnd/internal/apperrors"
	config2 "GoAPIBackEnd/internal/config"
	"GoAPIBackEnd/internal/database"
	"GoAPIBackEnd/internal/helpers"
	"GoAPIBackEnd/internal/type/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

var validate = validator.New()

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
	var dbConn = config2.Conn

	userBody := models.Users{}
	// Get user input
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

	// validate user inp
	if validationErr := validate.Struct(userBody); validationErr != nil {
		var errorMessages []string
		var validationErrors validator.ValidationErrors
		if errors.As(validationErr, &validationErrors) {
			for _, err := range validationErrors {
				errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.Tag()))
			}
		}
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "VALIDATION_ERROR", fmt.Errorf(strings.Join(errorMessages, "; "))))
		return
	}

	user, err := database.GetUser(dbConn, "", userBody.Name, "")
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}
	if len(user) > 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusConflict, "USER_EXISTS", fmt.Errorf("User already exists")))
		return
	}

	accessToke, refreshToken := helpers.GenerateToken(userBody.Email, userBody.Id, userBody.Role)
	nUser := models.Users{
		Id:           uuid.New().String(),
		Name:         userBody.Name,
		Password:     *helpers.HashPassword(userBody.Password),
		Email:        userBody.Email,
		Role:         userBody.Role,
		CoachId:      nil,
		Token:        &accessToke,
		RefreshToken: &refreshToken,
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
	var dbConn = config2.Conn

	userPayload := models.Users{}
	if err := ctx.ShouldBindJSON(&userPayload); err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	user, err := database.GetUser(dbConn, "", userPayload.Name, userPayload.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusUnauthorized, "INVALID_CREDENTIALS", fmt.Errorf("Invalid username or password")))
		return
	}
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}

	passwordIsValid, msg := helpers.VerifyPassword(user[0].Password, userPayload.Password)
	if !passwordIsValid {
		apperrors.GetAPIError(ctx, gin.H{"error": msg}, 0, apperrors.APIError{}.New(http.StatusUnauthorized, "INVALID_CREDENTIALS", fmt.Errorf("Invalid username or password")))
		return
	}

	token, refreshToken := helpers.GenerateToken(user[0].Email, user[0].Id, user[0].Role)
	helpers.UpdatedAllToken(token, refreshToken, user[0].Id)

	// Send it back
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", token, int(time.Now().Add(7*24*time.Hour).Unix()), "/", "", false, true)

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
