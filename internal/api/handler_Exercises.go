package api

import (
	"GoAPIBackEnd/internal/apperrors"
	"GoAPIBackEnd/internal/config"
	"GoAPIBackEnd/internal/database"
	"GoAPIBackEnd/internal/type/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// @Summary Create Exercise.
// @Tags Exercise
// @Produce json
// @Param request body models.Exercises true "Exercise"
// @Success 200 {object} models.APIError
// @Failure 208 {object} models.APIError
// @Failure 400 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /exercises [post]
func PostExercise(ctx *gin.Context) {
	var dbConn = config.Conn
	ex := models.Exercises{}
	err := ctx.ShouldBindJSON(&ex)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	if len(ex.Id) == 0 {
		uid, err := uuid.NewUUID()
		if err != nil {
			return
		}
		ex.Id = uid.String()
	}

	exercise, err := database.GetExercise(dbConn, ex.Id, "")
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}

	if len(exercise) > 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusConflict, "EXERCISE_EXISTS", fmt.Errorf("Exercise already exists")))
		return
	}

	err = database.PostExercise(dbConn, ex)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "INVALID_REQUEST", err))
		return
	}

	apperrors.GetAPIError(ctx, nil, http.StatusOK, nil)
}

// @Summary Get Exercise.
// @Tags Exercise
// @Produce json
// @Param id path string true "Exercise Id"
// @Success 200 {object} models.Exercises
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /exercises/{id} [get]
func GetExercise(ctx *gin.Context) {
	var dbConn = config.Conn
	id := ctx.Param("id")

	exercise, err := database.GetExercise(dbConn, id, "")
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}
	if len(exercise) == 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotFound, "EXERCISE_NOT_FOUND", fmt.Errorf("Exercise with id '%s' not found", id)))
		return
	}

	apperrors.GetAPIError(ctx, exercise[0], http.StatusOK, nil)
}

// @Summary	Query Exercises.
// @Tags		Exercise
// @Produce	json
// @Success	200		{array}		models.Exercises
// @Failure 404     {object}    models.APIError
// @Failure	500		{object}	models.APIError
// @Router		/exercises [get]
func QueryExercises(ctx *gin.Context) {
	var dbConn = config.Conn

	exercise, err := database.GetExercise(dbConn, "", "")
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}
	if len(exercise) == 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotFound, "EXERCISE_NOT_FOUND", fmt.Errorf("Exercise not found")))
		return
	}

	apperrors.GetAPIError(ctx, exercise, http.StatusOK, nil)
}
