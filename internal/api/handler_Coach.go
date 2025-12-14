package api

import (
	"GoAPIBackEnd/internal/apperrors"
	"GoAPIBackEnd/internal/database"
	"GoAPIBackEnd/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// @Summary Create Coach.
// @Tags Coach
// @Produce json
// @Param request body model.Coach true "Coach"
// @Success 200 {object} model.APIError
// @Failure 208 {object} model.APIError
// @Failure 400 {object} model.APIError
// @Failure 500 {object} model.APIError
// @Router /coach [post]
func PostCoach(ctx *gin.Context) {
	var dbConn = database.Conn
	co := models.Coach{}
	err := ctx.ShouldBindJSON(&co)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	if len(co.Id) == 0 {
		uid, err := uuid.NewUUID()
		if err != nil {
			return
		}
		co.Id = uid.String()
	}

	coach, err := database.GetCoach(dbConn, co.Name, false)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}

	if len(coach) > 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusConflict, "EXERCISE_EXISTS", fmt.Errorf("Exercise already exists")))
		return
	}

	err = database.PostCoach(dbConn, co)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "INVALID_REQUEST", err))
		return
	}

	apperrors.GetAPIError(ctx, nil, http.StatusOK, nil)
}

// @Summary Get Coach.
// @Tags Coach
// @Produce json
// @Param Name path string true "Coach Name"
// @Success 200 {object} model.Exercise
// @Failure 404 {object} model.APIError
// @Failure 500 {object} model.APIError
// @Router /Coach/{name} [get]
func GetCoach(ctx *gin.Context) {
	var dbConn = database.Conn
	name := ctx.Param("name")

	coach, err := database.GetCoach(dbConn, name, true)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}
	if len(coach) == 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotFound, "EXERCISE_NOT_FOUND", fmt.Errorf("Exercise with name '%s' not found", name)))
		return
	}

	apperrors.GetAPIError(ctx, coach, http.StatusOK, nil)
}
