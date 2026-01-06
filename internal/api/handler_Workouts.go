package api

import (
	"GoAPIBackEnd/internal/apperrors"
	"GoAPIBackEnd/internal/config"
	"GoAPIBackEnd/internal/database"
	"GoAPIBackEnd/internal/type/misc"
	"GoAPIBackEnd/internal/type/models"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var dbMutex sync.Mutex

const pathCSV = "C:\\Users\\User\\GolandProjects\\GoAPIBackEnd\\excel\\"

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
	var dbConn = config.Conn
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

	wk, err := database.GetWorkout(dbConn, workout.Id, "", "", false)
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
	var dbConn = config.Conn
	id := ctx.Param("id")

	workout, err := database.GetWorkout(dbConn, id, "", "", true)
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
// @Description Query workouts with optional user filtering
// @Tags		Workouts
// @Produce	json
// @Param       userId  query    string  false  "Filter by user Id"
// @Success	200		{array}		models.Workouts
// @Failure 404     {object}    models.APIError
// @Failure	500		{object}	models.APIError
// @Router		/workouts [get]
func QueryWorkouts(ctx *gin.Context) {
	var dbConn = config.Conn

	// Get all query parameters
	userId := ctx.Query("userId")

	workout, err := database.GetWorkout(dbConn, "", "", userId, true)
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
	var dbConn = config.Conn
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

	w, err := database.GetWorkout(dbConn, workout.Id, "", "", false)
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
	var dbConn = config.Conn
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
	var dbConn = config.Conn
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

// @Summary	Load Workouts
// @Tags		Workouts
// @Produce	json
// @Param		request	body		misc.LoadFiles	true	"Load Files"
// @Success	200		{object}	models.LoadWorkouts
// @Failure	400		{object}	models.APIError
// @Failure	408     {object}	models.APIError
// @Failure	500		{object}	models.APIError
// @Router		/load/workouts [post]
func LoadExcels(ctx *gin.Context) {
	files := misc.LoadFiles{}
	err := ctx.ShouldBindJSON(&files)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	resultCh := make(chan models.LoadWorkouts, len(files.Workouts))
	var wg sync.WaitGroup

	for _, fileName := range files.Workouts {
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()
			resultCh <- loadCSVFile(fname)
		}(fileName)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	results := make([]models.LoadWorkouts, 0, len(files.Workouts))
	timeout := time.After(30 * time.Second)

	for {
		select {
		case res, ok := <-resultCh:
			if !ok {
				apperrors.GetAPIError(ctx, gin.H{"results": results}, http.StatusOK, nil)
				return
			}
			results = append(results, res)

		case <-timeout:
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"error":   "Timeout waiting for load",
				"results": results,
			})
			return
		}
	}

	apperrors.GetAPIError(ctx, gin.H{"results": results}, http.StatusOK, nil)
}

func loadCSVFile(fileName string) models.LoadWorkouts {
	path := pathCSV + fileName
	file, err := os.Open(path)
	if err != nil {
		return models.LoadWorkouts{Excel: fileName, Success: false, Error: err}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read() // Skip header
	records, err := reader.ReadAll()
	if err != nil {
		return models.LoadWorkouts{Excel: fileName, Success: false, Error: err}
	}

	dbConn := config.Conn

	inserted := 0
	failed := 0

	for _, record := range records {
		notes := record[3]
		workout := models.Workouts{
			Id:        uuid.NewString(),
			UserId:    strings.TrimSpace(record[1]),
			Name:      strings.TrimSpace(record[2]),
			Notes:     &notes,
			CreatedAt: time.Now(),
		}

		dbMutex.Lock()
		err := database.PostWorkout(dbConn, workout)
		dbMutex.Unlock()

		if err != nil {
			fmt.Printf("FAILED: %v\n", err)
			failed++
		} else {
			inserted++
		}
	}

	return models.LoadWorkouts{
		Excel:    fileName,
		Success:  failed == 0,
		Inserted: inserted,
		Failed:   failed,
		Error:    nil,
	}
}
