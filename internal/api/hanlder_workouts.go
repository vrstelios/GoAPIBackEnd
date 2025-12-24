package api

import (
	"GoAPIBackEnd/internal/apperrors"
	"GoAPIBackEnd/internal/database"
	"GoAPIBackEnd/internal/type/misc"
	"GoAPIBackEnd/internal/type/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"net/http"
)

// @Summary Create Workout.
// @Tags Workouts
// @Produce json
// @Param request body models.Workouts true "Workout"
// @Success 200 {object} models.APIError
// @Failure 208 {object} models.APIError
// @Failure 400 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /workouts [post]
func PostWorkout(ctx *gin.Context) {
	var dbConn = database.Conn
	workout := models.Workouts{}
	err := ctx.ShouldBindJSON(&workout)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	if len(workout.Id) == 0 {
		uid, err := uuid.NewUUID()
		if err != nil {
			return
		}
		workout.Id = uid.String()
	}

	wk, err := database.GetWorkout(dbConn, workout.Id, "", false)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}

	if len(wk) > 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusConflict, "EXERCISE_EXISTS", fmt.Errorf("Workout already exists")))
		return
	}

	err = database.PostWorkout(dbConn, workout)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "INVALID_REQUEST", err))
		return
	}

	apperrors.GetAPIError(ctx, nil, http.StatusOK, nil)
}

// @Summary Get Workout.
// @Tags Workouts
// @Produce json
// @Param id path string true "Workouts Id"
// @Success 200 {object} models.Workouts
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /workouts/{id} [get]
func GetWorkout(ctx *gin.Context) {
	var dbConn = database.Conn
	id := ctx.Param("id")

	workout, err := database.GetWorkout(dbConn, id, "", true)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}
	if len(workout) == 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotFound, "EXERCISE_NOT_FOUND", fmt.Errorf("Workout with id '%s' not found", id)))
		return
	}

	apperrors.GetAPIError(ctx, workout[0], http.StatusOK, nil)
}

// @Summary	Query Workouts.
// @Tags		Workouts
// @Produce	json
// @Success	200		{array}		models.Workouts
// @Failure 404     {object}    models.APIError
// @Failure	500		{object}	models.APIError
// @Router		/workouts [get]
func QueryWorkouts(ctx *gin.Context) {
	var dbConn = database.Conn

	workout, err := database.GetWorkout(dbConn, "", "", true)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}
	if len(workout) == 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotFound, "EXERCISE_NOT_FOUND", fmt.Errorf("Exercise not found")))
		return
	}

	apperrors.GetAPIError(ctx, workout, http.StatusOK, nil)
}

// @Summary Update Workout.
// @Tags Workouts
// @Produce json
// @Param id path string false "Workout"
// @Param request body models.Workouts true "Workouts"
// @Success 200 {array}	 models.Workouts
// @Failure 400 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /workouts/{id} [put]
func PutWorkout(ctx *gin.Context) {
	var dbConn = database.Conn
	workout := models.Workouts{}
	err := ctx.ShouldBindJSON(&workout)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	id := ctx.Param("id")
	if workout.Id != "" && workout.Id != id {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "ID_MISMATCH", fmt.Errorf("Id in body must match URL")))
		return
	}

	err = database.PutWorkout(dbConn, workout)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "INVALID_REQUEST", err))
		return
	}

	w, err := database.GetWorkout(dbConn, workout.Id, "", false)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}

	if len(w) == 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusConflict, "EXERCISE_EXISTS", fmt.Errorf("Workout don't found")))
		return
	}

	apperrors.GetAPIError(ctx, w[0], http.StatusOK, nil)
}

// @Summary Delete Workout.
// @Tags Workouts
// @Produce json
// @Param id path string false "Workout"
// @Success 200 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /workouts/{id} [delete]
func DelWorkout(ctx *gin.Context) {
	var dbConn = database.Conn
	id := ctx.Param("id")

	err := database.DelWorkout(dbConn, id)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotFound, "EXERCISE_NOT_FOUND", fmt.Errorf("Workout with id '%s' not found", id)))
		return
	}

	apperrors.GetAPIError(ctx, nil, http.StatusNoContent, nil)
}

// @Summary	Query Workouts (server side filtering & paging)
// @Tags		Workouts
// @Produce	json
// @Param		request	body		misc.QueryWorkoutsRequest	true	"QueryWorkoutsRequest"
// @Success	200		{object}	misc.QueryWorkoutsResponse
// @Failure	400		{object}	models.APIError
// @Failure	500		{object}	models.APIError
// @Router		/workouts/query [post]
func QueryWorkoutsV2(ctx *gin.Context) {
	var dbConn = database.Conn
	queryRequest := misc.QueryWorkoutsRequest{}
	err := ctx.ShouldBindJSON(&queryRequest)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	var filtered []models.Workouts
	orderStr := "ASC"
	if queryRequest.Ordering == 0 {
		orderStr = "DESC"
	}

	sql := fmt.Sprintf(`
	SELECT id, user_id, name, notes, scheduled_at, created_at
	FROM workouts
	WHERE ($1 = '' OR name ILIKE '%%' || $1 || '%%')
	ORDER BY name %s
	LIMIT $2 OFFSET $3`, orderStr)

	offset := (queryRequest.Paging.Page - 1) * queryRequest.Paging.PageSize
	limit := queryRequest.Paging.PageSize

	rows, err := dbConn.Query(context.Background(), sql, queryRequest.Search, limit, offset)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var w models.Workouts
		err := rows.Scan(&w.Id, &w.UserId, &w.Name, &w.Notes, &w.ScheduledAt, &w.CreatedAt)
		if err != nil {
			apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
			return
		}
		filtered = append(filtered, w)
	}

	totalRows := len(filtered)
	if totalRows == 0 {
		apperrors.GetAPIError(ctx, misc.QueryWorkoutsResponse{TotalResults: totalRows, TotalPages: 0, Workouts: []models.Workouts{}}, http.StatusOK, nil)
		return
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(queryRequest.Paging.PageSize)))
	response := misc.QueryWorkoutsResponse{
		TotalResults: totalRows,
		TotalPages:   totalPages,
		Workouts:     filtered,
	}
	apperrors.GetAPIError(ctx, response, http.StatusOK, nil)
}
